package main

import (
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
	file, err := initLogTxt()
	if err != nil {
		log.Println(err)
		return
	}

	defer onBeforeExit(file) // called at the end if no SIGINT or SIGTERM is received
	go listenToSIGINTAndSIGTERM(func() { onBeforeExit(file) })

	// init game state
	GAME_CONFIG := createGameConfig()
	canvas := createCanvas(GAME_CONFIG)
	snakeHeadPos := Vector2{0, 0}
	snakeDirection := Vector2{1, 0}

	// game loop
	for { // "'while' is spelled 'for' in Go"

		// TODO : poll input

		// game logic - OnUpdate

		// update snake head pos
		snakeHeadPos = snakeHeadPos.Add(snakeDirection)

		// wrap around canvas edge
		snakeHeadPos.x = snakeHeadPos.x % GAME_CONFIG.CANVAS_SIZE.x
		snakeHeadPos.y = snakeHeadPos.y % GAME_CONFIG.CANVAS_SIZE.y
		if snakeHeadPos.x < 0 {
			snakeHeadPos.x += GAME_CONFIG.CANVAS_SIZE.x
		}
		if snakeHeadPos.y < 0 {
			snakeHeadPos.y += GAME_CONFIG.CANVAS_SIZE.y
		}

		// game logic - OnAfterUpdate
		// TODO : reset input buffer

		// draw
		resetCanvas(canvas, GAME_CONFIG)
		drawChar(canvas, snakeHeadPos, GAME_CONFIG.SNAKE_CHAR, GAME_CONFIG)

		// render
		buffer := canvasToStringBuffer(canvas)
		clearAndDrawBuffer(buffer)

		// frame
		// TODO : see multiplayer book suggested architecture
		time.Sleep(time.Duration(GAME_CONFIG.DELTA_TIME_MS) * time.Millisecond)
	}
}

func listenToSIGINTAndSIGTERM(onBeforeExit func()) {
	// create a channel, type os.Signal, buffer 1 (required by signal.Notify)
	channel := make(chan os.Signal, 1)

	// notify channel on os.Interrupt (SIGINT) or SIGTERM
	signal.Notify(channel, os.Interrupt, syscall.SIGTERM)

	// blocked until receives a notification from channel
	receivedSignal := <-channel // go does not allow unused variables
	log.Println("receivedSignal:", receivedSignal)

	// need to manually call cuz deferred calls will not be called after exit
	onBeforeExit()

	// caught the signal myself, so also need to exit myself, as it overrides the default behavior
	os.Exit(0)
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

func createCanvas(config GameConfig) [][]string {
	canvas := [][]string{}

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

func resetCanvas(canvas [][]string, config GameConfig) [][]string {
	for y := 0; y < config.CANVAS_SIZE.y; y++ {
		for x := 0; x < config.CANVAS_SIZE.x; x++ {
			drawChar(canvas, Vector2{x, y}, config.BG_CHAR, config)
		}
	}
	return canvas
}

func drawChar(canvas [][]string, position Vector2, char string, config GameConfig) {
	// canvas is drawn from top to bottom but game coordinates is from bottom to top
	canvas[config.TERM_SIZE.y-(config.PADDING.y+position.y)-2][config.PADDING.x+position.x] = char
}

func canvasToStringBuffer(canvas [][]string) string {
	var buffer strings.Builder

	for _, row := range canvas {
		for _, element := range row {
			buffer.WriteString(element)
		}
		buffer.WriteString("\n")
	}

	return buffer.String()
}
func clearAndDrawBuffer(buffer string) {
	ansi.ClearAndDrawBuffer(buffer)
}

func onBeforeExit(logFile *os.File) {
	log.Println("onBeforeExit: cleanup")
	ansi.Cleanup()

	log.Println("onBeforeExit: close log.txt")
	logFile.Close()
}
