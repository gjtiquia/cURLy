package main

import (
	"log"
	"math"
	"time"
)

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

func createGameConfig(termWidth, termHeight int) GameConfig {
	const FPS = 10
	const DELTA_TIME_MS = 1000 / FPS

	CANVAS_SIZE := Vector2{40, 10}
	BORDER_THICKNESS := Vector2{1, 1}

	// TODO : should make effort to support resize...? right now assumes static
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
