package main

import (
	"log"
	"math"
	"os"
	"strings"
	"time"

	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
	"github.com/pkg/errors"
)

// TODO : LMAO SHOWS NOTHING

func main() {
	// logging setup
	err, logPanicAndCloseFile := InitLogFile("log.txt")
	if err != nil {
		log.Panicf("%+v", errors.WithStack(err))
	}
	defer logPanicAndCloseFile()

	// tcell setup
	s, _, err, finalizeScreen := InitTCellScreen()
	if err != nil {
		log.Panicf("%+v", err)
	}
	defer finalizeScreen()

	// gameSetup
	gameConfig, gameState, canvas, inputBuffer := createGame(s)

	for {

		startTime := time.Now()

		inputBuffer, isExit := DrainTCellEvents(s, inputBuffer)
		if isExit {
			return
		}

		runGameLoop(gameConfig, gameState, canvas, inputBuffer, s)

		elapsedTime := time.Since(startTime)
		remainingTime := gameConfig.DELTA_TIME - elapsedTime
		if remainingTime > 0 {
			time.Sleep(remainingTime)
		}
	}

}

func InitLogFile(filename string) (err error, logPanicAndCloseFile func()) {
	// truncate means delete contents on open, create if doesnt exist, write-only
	const fileFlags = os.O_TRUNC | os.O_CREATE | os.O_WRONLY

	// read = 4, write = 2, execute = 1; 6 = 4 + 2 (read write); 0 = octal; 666 = owner/group/others
	const filePerm = 0666

	file, err := os.OpenFile(filename, fileFlags, filePerm)
	if err != nil {
		return err, nil
	}

	log.SetOutput(file)
	logPanicAndCloseFile = func() {
		defer file.Close()
		if r := recover(); r != nil {
			log.Println("logging panic before file close")
			log.Panicf("%+v", r)
		}
	}
	return nil, logPanicAndCloseFile
}

func InitTCellScreen() (s tcell.Screen, defStyle tcell.Style, err error, finalizeScreen func()) {
	defStyle = tcell.StyleDefault.Background(color.Reset).Foreground(color.Reset)

	s, err = tcell.NewScreen()
	if err != nil {
		return nil, defStyle, err, nil
	}
	if err = s.Init(); err != nil {
		return nil, defStyle, err, nil
	}

	// Set default text style
	s.SetStyle(defStyle)

	// Clear screen
	s.Clear()

	finalizeScreen = func() {
		// You have to catch panics in a defer, clean up, and re-raise them - otherwise your application can die without leaving any diagnostic trace.
		// https://github.com/gdamore/tcell/blob/main/TUTORIAL.md
		maybePanic := recover()
		s.Fini()
		if maybePanic != nil {
			panic(maybePanic)
		}
	}

	return s, defStyle, nil, finalizeScreen
}

func DrainTCellEvents(s tcell.Screen, inputBuffer []InputAction) ([]InputAction, bool) {
	inputBuffer = inputBuffer[:0]

	for {
		// Update screen
		s.Show()

		// Poll event (can be used in select statement as well)
		ev := <-s.EventQ()

		// Process event
		switch ev := ev.(type) {

		case *tcell.EventResize:
			s.Sync()

		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
				return inputBuffer, true
			}

			switch key, str := ev.Key(), ev.Str(); {
			case key == tcell.KeyUp, str == "w":
				inputBuffer = append(inputBuffer, Up)
			case key == tcell.KeyDown, str == "s":
				inputBuffer = append(inputBuffer, Up)
			case key == tcell.KeyLeft, str == "a":
				inputBuffer = append(inputBuffer, Up)
			case key == tcell.KeyRight, str == "d":
				inputBuffer = append(inputBuffer, Up)
			}

		default:
			return inputBuffer, false
		}
	}
}

func createGame(s tcell.Screen) (GameConfig, *GameState, GameCanvas, []InputAction) {
	gameConfig := createGameConfig(s)
	gameState := createGameState()
	canvas := createCanvas(gameConfig)
	inputBuffer := make([]InputAction, 0)
	return gameConfig, gameState, canvas, inputBuffer
}

func runGameLoop(gameConfig GameConfig, gameState *GameState, canvas GameCanvas, inputBuffer []InputAction, s tcell.Screen) {
	// game logic
	gameState.onUpdate(gameConfig, inputBuffer)

	// draw
	canvas.resetCanvas(gameConfig)
	gameState.onDraw(gameConfig, canvas)

	// render
	buffer := canvas.toStringBuffer()
	s.PutStr(0, 0, buffer)

	// TODO : rethink on whether should PutStr of ENTIRE CANVAS, or that we should actually abstract that out...
	// cuz, tcell can simply replace one char at a time, so, that "might" be more efficient?
	// ALTHOUGH we are porting to the web lol, so maybe
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

type Vector2 struct {
	x int
	y int
}

func (this Vector2) Add(other Vector2) Vector2 {
	return Vector2{x: this.x + other.x, y: this.y + other.y}
}

type GameConfig struct {
	FPS        int
	DELTA_TIME time.Duration

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

func createGameConfig(s tcell.Screen) GameConfig {
	const FPS = 10
	const DELTA_TIME_MS = 1000 / FPS

	CANVAS_SIZE := Vector2{40, 10}
	BORDER_THICKNESS := Vector2{1, 1}

	// TODO : should make effort to support resize...? right now assumes static
	termWidth, termHeight := s.Size()

	paddingX := int(math.Floor(float64(termWidth-CANVAS_SIZE.x) / 2))
	paddingY := int(math.Floor(float64(termHeight-CANVAS_SIZE.y) / 2))

	log.Printf("term size: %vx%v", termWidth, termHeight)
	log.Printf("padding: %vx%v", paddingX, paddingY)
	log.Printf("canvas: %vx%v", CANVAS_SIZE.x, CANVAS_SIZE.y)

	return GameConfig{
		FPS:        FPS,
		DELTA_TIME: 1000 / FPS * time.Millisecond,

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

func createGameState() *GameState {
	gameState := GameState{
		snakeHeadPos:   Vector2{0, 0},
		snakeDirection: Vector2{1, 0},
	}
	return &gameState
}

func (this *GameState) onUpdate(gameConfig GameConfig, inputBuffer []InputAction) {

	// TODO : for now just get the most recent input action
	if len(inputBuffer) > 0 {
		inputAction := inputBuffer[len(inputBuffer)-1]

		switch {
		case inputAction == Up:
			this.snakeDirection = Vector2{0, 1}
		case inputAction == Down:
			this.snakeDirection = Vector2{0, -1}
		case inputAction == Left:
			this.snakeDirection = Vector2{-1, 0}
		case inputAction == Right:
			this.snakeDirection = Vector2{1, 0}
		}
	}

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

	// Define regions (all coordinates are in terminal space, 0-indexed from top-left)
	// Border surrounds the canvas, padding surrounds the border
	borderLeft := config.PADDING.x - config.BORDER_THICKNESS.x
	borderRight := config.PADDING.x + config.CANVAS_SIZE.x // exclusive for canvas, inclusive for border
	borderTop := config.PADDING.y - config.BORDER_THICKNESS.y
	borderBottom := config.PADDING.y + config.CANVAS_SIZE.y // exclusive for canvas, inclusive for border

	canvasLeft := config.PADDING.x
	canvasRight := config.PADDING.x + config.CANVAS_SIZE.x // exclusive
	canvasTop := config.PADDING.y
	canvasBottom := config.PADDING.y + config.CANVAS_SIZE.y // exclusive

	for y := 0; y < config.TERM_SIZE.y; y++ {
		row := []string{}
		for x := 0; x < config.TERM_SIZE.x; x++ {
			// Check if in canvas region
			if x >= canvasLeft && x < canvasRight && y >= canvasTop && y < canvasBottom {
				row = append(row, config.BG_CHAR)
				continue
			}

			// Check if in border region (horizontal borders - top and bottom)
			if y >= borderTop && y < canvasTop && x >= borderLeft && x < borderRight+config.BORDER_THICKNESS.x {
				row = append(row, config.BORDER_X_CHAR)
				continue
			}
			if y >= canvasBottom && y < borderBottom+config.BORDER_THICKNESS.y && x >= borderLeft && x < borderRight+config.BORDER_THICKNESS.x {
				row = append(row, config.BORDER_X_CHAR)
				continue
			}

			// Check if in border region (vertical borders - left and right)
			if y >= canvasTop && y < canvasBottom {
				if x >= borderLeft && x < canvasLeft {
					row = append(row, config.BORDER_Y_CHAR)
					continue
				}
				if x >= canvasRight && x < canvasRight+config.BORDER_THICKNESS.x {
					row = append(row, config.BORDER_Y_CHAR)
					continue
				}
			}

			// Everything else is padding
			row = append(row, config.PADDING_CHAR)
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
	// game y=0 is at bottom of canvas, which maps to terminal row: PADDING.y + CANVAS_SIZE.y - 1
	// game y=CANVAS_SIZE.y-1 is at top, which maps to terminal row: PADDING.y
	termY := config.PADDING.y + (config.CANVAS_SIZE.y - 1 - position.y)
	termX := config.PADDING.x + position.x
	this[termY][termX] = char
}

func (canvas GameCanvas) toStringBuffer() string {
	var buffer strings.Builder

	for _, row := range canvas {
		for _, element := range row {
			buffer.WriteString(element)
		}
		// Raw mode: terminal does not translate \n to \r\n
		// so LF (\n) alone only moves down, not to column 0
		// Use \r\n so each line starts at the left
		buffer.WriteString("\r\n")
	}

	return buffer.String()
}
