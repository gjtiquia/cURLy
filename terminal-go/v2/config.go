package main

import (
	"log"
	"math"
	"time"

	"github.com/gjtiquia/cURLy/terminal-go/v2/internals/vector2"
)

type GameConfig struct {
	FPS        int
	DELTA_TIME time.Duration

	TERM_SIZE        vector2.Type
	BORDER_THICKNESS vector2.Type
	CANVAS_SIZE      vector2.Type
	PADDING          vector2.Type

	PADDING_CHAR  string
	BORDER_X_CHAR string
	BORDER_Y_CHAR string
	BG_CHAR       string
	SNAKE_CHAR    string
	FOOD_CHAR     string
}

func createGameConfig(termWidth, termHeight int) GameConfig {
	const FPS = 10
	const DELTA_TIME_MS = 1000 / FPS

	CANVAS_SIZE := vector2.New(40, 10)
	BORDER_THICKNESS := vector2.New(1, 1)

	// TODO : should make effort to support resize...? right now assumes static
	paddingX := int(math.Floor(float64(termWidth-CANVAS_SIZE.X) / 2))
	paddingY := int(math.Floor(float64(termHeight-CANVAS_SIZE.Y) / 2))

	log.Printf("term size: %vx%v", termWidth, termHeight)
	log.Printf("padding: %vx%v", paddingX, paddingY)
	log.Printf("canvas: %vx%v", CANVAS_SIZE.X, CANVAS_SIZE.Y)

	return GameConfig{
		FPS:        FPS,
		DELTA_TIME: 1000 / FPS * time.Millisecond,

		CANVAS_SIZE:      CANVAS_SIZE,
		BORDER_THICKNESS: BORDER_THICKNESS,
		TERM_SIZE:        vector2.New(termWidth, termHeight),
		PADDING:          vector2.New(paddingX, paddingY),

		PADDING_CHAR:  " ",
		BORDER_X_CHAR: "-",
		BORDER_Y_CHAR: "|",
		BG_CHAR:       " ",
		SNAKE_CHAR:    "O",
		FOOD_CHAR:     "*",
	}
}
