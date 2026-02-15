package main

import (
	"errors"
	"log"
	"math"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gjtiquia/cURLy/terminal-go/ansi"
	"golang.org/x/term"
)

func main() {
	// logging setup
	file, err := initLogTxt()
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()
	defer log.Println("main: terminating...")

	// Ctrl-C setup
	errCh := make(chan error)
	go listenForSIGINTAndSIGTERM(errCh) // manually handle Ctrl-C or else defers wont call

	// display setup
	ansi.ClearAndHideCursor()
	defer ansi.ClearAndShowCursor()

	// input setup
	// oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }
	// defer term.Restore(int(os.Stdin.Fd()), oldState)

	inputCh := make(chan InputAction)
	go listenForInput(inputCh, errCh)

	// game setup
	gameConfig, gameState, canvas := createGame()

	for {
		select {

		case error := <-errCh:
			log.Println("main.errCh:", error)
			return

		case inputAction := <-inputCh:
			if inputAction == Exit {
				errCh <- errors.New("input: Exit")
				continue
			}
			handleInput(inputAction)

		default:
			// TODO : see multiplayer book suggested architecture
			runGameLoop(gameConfig, gameState, canvas)
			time.Sleep(time.Duration(gameConfig.DELTA_TIME_MS) * time.Millisecond)
		}
	}
}

func createGame() (GameConfig, GameState, GameCanvas) {
	gameConfig := createGameConfig()
	gameState := createGameState()
	canvas := createCanvas(gameConfig)
	return gameConfig, gameState, canvas
}

func runGameLoop(gameConfig GameConfig, gameState GameState, canvas GameCanvas) {
	// game logic
	gameState.onUpdate(gameConfig)

	// draw
	canvas.resetCanvas(gameConfig)
	gameState.onDraw(gameConfig, canvas)

	// render
	buffer := canvas.toStringBuffer()
	ansi.ClearAndDrawBuffer(buffer)
}

func listenForSIGINTAndSIGTERM(errCh chan error) {
	// create a channel, type os.Signal, buffer 1 (required by signal.Notify)
	channel := make(chan os.Signal, 1)

	// notify channel on os.Interrupt (SIGINT) or SIGTERM
	signal.Notify(channel, os.Interrupt, syscall.SIGTERM)

	// blocked until receives a notification from channel
	receivedSignal := <-channel

	// notify outside
	errCh <- errors.New("listenForSIGINTAndSIGTERM: " + receivedSignal.String())
}

type InputAction int

const (
	None InputAction = iota
	Up
	Down
	Left
	Right
	Exit
)

func listenForInput(inputCh chan InputAction, errCh chan error) {
	buf := make([]byte, 3) // read up to 3 bytes, eg. arrows = 3 byte escape sequence
	for {
		n, err := os.Stdin.Read(buf)
		if err != nil {
			errCh <- errors.New("listenForInput: " + err.Error())
			return
		}

		// Handle Escape Sequences (Arrows)
		if n == 3 && buf[0] == 27 && buf[1] == 91 {
			switch buf[2] {
			case 65:
				inputCh <- Up
			case 66:
				inputCh <- Down
			case 67:
				inputCh <- Right
			case 68:
				inputCh <- Left
			}
			continue
		}

		// Handle Single Bytes
		switch buf[0] {
		case 3: // Ctrl+C
			log.Println("listenForInput: Ctrl+C")
			inputCh <- Exit
		case 'w':
			inputCh <- Up
		case 's':
			inputCh <- Down
		case 'a':
			inputCh <- Left
		case 'd':
			inputCh <- Right
		}
	}
}

func handleInput(action InputAction) {
	log.Println("input", action)
}

func initLogTxt() (*os.File, error) {
	// truncate means delete contents on open, create if doesnt exist, write-only
	const fileFlags = os.O_TRUNC | os.O_CREATE | os.O_WRONLY

	// read = 4, write = 2, execute = 1; 6 = 4 + 2 (read write); 0 = octal; 666 = owner/group/others
	const filePerm = 0666

	file, err := os.OpenFile("log.txt", fileFlags, filePerm)
	if err != nil {
		return nil, err
	}

	log.SetOutput(file)
	return file, nil
}

type Vector2 struct {
	x int
	y int
}

func (this Vector2) Add(other Vector2) Vector2 {
	return Vector2{x: this.x + other.x, y: this.y + other.y}
}

type GameConfig struct {
	FPS           int
	DELTA_TIME_MS int

	TERM_SIZE        Vector2
	BORDER_THICKNESS Vector2
	CANVAS_SIZE      Vector2
	PADDING          Vector2

	PADDING_CHAR  string
	BORDER_X_CHAR string
	BORDER_Y_CHAR string
	BG_CHAR       string
	SNAKE_CHAR    string
}

func createGameConfig() GameConfig {
	const FPS = 10
	const DELTA_TIME_MS = 1000 / FPS

	CANVAS_SIZE := Vector2{40, 10}
	BORDER_THICKNESS := Vector2{1, 1}

	currentTermFd := int(os.Stdout.Fd())
	termWidth, termHeight, err := term.GetSize(currentTermFd)
	if err != nil {
		termWidth = CANVAS_SIZE.x + BORDER_THICKNESS.x
		termHeight = CANVAS_SIZE.y + BORDER_THICKNESS.y
		log.Println("Unable to get term size! using fallbacks...")
	}

	paddingX := int(math.Floor(float64(termWidth-CANVAS_SIZE.x) / 2))
	paddingY := int(math.Floor(float64(termHeight-CANVAS_SIZE.y) / 2))

	log.Printf("term size: %vx%v", termWidth, termHeight)
	log.Printf("padding: %vx%v", paddingX, paddingY)
	log.Printf("canvas: %vx%v", CANVAS_SIZE.x, CANVAS_SIZE.y)

	return GameConfig{
		FPS:           10,
		DELTA_TIME_MS: 1000 / FPS,

		CANVAS_SIZE:      CANVAS_SIZE,
		BORDER_THICKNESS: BORDER_THICKNESS,
		TERM_SIZE:        Vector2{termWidth, termHeight},
		PADDING:          Vector2{paddingX, paddingY},

		PADDING_CHAR:  " ",
		BORDER_X_CHAR: "-",
		BORDER_Y_CHAR: "|",
		BG_CHAR:       " ",
		SNAKE_CHAR:    "x",
	}
}

type GameState struct {
	snakeHeadPos   Vector2
	snakeDirection Vector2
}

func createGameState() GameState {
	return GameState{
		snakeHeadPos:   Vector2{0, 0},
		snakeDirection: Vector2{1, 0},
	}
}

func (this GameState) onUpdate(gameConfig GameConfig) {
	// update snake head pos
	this.snakeHeadPos = this.snakeHeadPos.Add(this.snakeDirection)

	// wrap around canvas edge
	this.snakeHeadPos.x = this.snakeHeadPos.x % gameConfig.CANVAS_SIZE.x
	this.snakeHeadPos.y = this.snakeHeadPos.y % gameConfig.CANVAS_SIZE.y
	if this.snakeHeadPos.x < 0 {
		this.snakeHeadPos.x += gameConfig.CANVAS_SIZE.x
	}
	if this.snakeHeadPos.y < 0 {
		this.snakeHeadPos.y += gameConfig.CANVAS_SIZE.y
	}
}

func (this GameState) onDraw(gameConfig GameConfig, canvas GameCanvas) {
	canvas.drawChar(this.snakeHeadPos, gameConfig.SNAKE_CHAR, gameConfig)
}

type GameCanvas [][]string

func createCanvas(config GameConfig) GameCanvas {
	canvas := GameCanvas{}

	// upper padding
	for y := 0; y < config.TERM_SIZE.y; y++ {
		row := []string{}
		for x := 0; x < config.TERM_SIZE.x; x++ {
			switch {

			case y < config.PADDING.y-config.BORDER_THICKNESS.y:
				row = append(row, config.PADDING_CHAR)

			case y < config.PADDING.y:
				switch {
				case x < config.PADDING.x-config.BORDER_THICKNESS.x:
					row = append(row, config.PADDING_CHAR)
				case x > config.PADDING.x+config.CANVAS_SIZE.x:
					row = append(row, config.PADDING_CHAR)
				default:
					row = append(row, config.BORDER_X_CHAR)
				}

			case y < config.PADDING.y+config.CANVAS_SIZE.y:
				switch {

				case x < config.PADDING.x-config.BORDER_THICKNESS.x:
					row = append(row, config.PADDING_CHAR)
				case x < config.PADDING.x:
					row = append(row, config.BORDER_Y_CHAR)

				case x > config.PADDING.x+config.CANVAS_SIZE.x:
					row = append(row, config.PADDING_CHAR)
				case x > config.PADDING.x+config.CANVAS_SIZE.x-config.BORDER_THICKNESS.x:
					row = append(row, config.BORDER_Y_CHAR)

				default:
					row = append(row, config.BG_CHAR)
				}

			case y >= config.PADDING.y+config.CANVAS_SIZE.y+config.BORDER_THICKNESS.y:
				row = append(row, config.PADDING_CHAR)

			default:
				switch {
				case x < config.PADDING.x-config.BORDER_THICKNESS.x:
					row = append(row, config.PADDING_CHAR)
				case x > config.PADDING.x+config.CANVAS_SIZE.x:
					row = append(row, config.PADDING_CHAR)
				default:
					row = append(row, config.BORDER_X_CHAR)
				}
			}
		}
		canvas = append(canvas, row)
	}

	return canvas
}

func (this GameCanvas) resetCanvas(config GameConfig) {
	for y := 0; y < config.CANVAS_SIZE.y; y++ {
		for x := 0; x < config.CANVAS_SIZE.x; x++ {
			this.drawChar(Vector2{x, y}, config.BG_CHAR, config)
		}
	}
}

func (this GameCanvas) drawChar(position Vector2, char string, config GameConfig) {
	// canvas is drawn from top to bottom but game coordinates is from bottom to top
	this[config.TERM_SIZE.y-(config.PADDING.y+position.y)-2][config.PADDING.x+position.x] = char
}

func (canvas GameCanvas) toStringBuffer() string {
	var buffer strings.Builder

	for _, row := range canvas {
		for _, element := range row {
			buffer.WriteString(element)
		}
		buffer.WriteString("\n")
	}

	return buffer.String()
}
