package game

import (
	"log"
	"math"
	"time"

	"github.com/gjtiquia/cURLy/internal/vector2"
)

type GameConfig struct {
	FPS        int
	DELTA_TIME time.Duration

	TERM_SIZE        vector2.Type
	BORDER_THICKNESS vector2.Type
	CANVAS_SIZE      vector2.Type
	PADDING          vector2.Type

	TITLE  string
	FOOTER string
}

func createGameConfig(termSize vector2.Type) GameConfig {
	FPS := 8

	CANVAS_SIZE := vector2.New(20, 8)
	BORDER_THICKNESS := vector2.New(1, 1)

	// uncomment to debug win case
	// FPS = 5
	// CANVAS_SIZE = vector2.New(4, 2)

	// TODO : should make effort to support resize...? right now assumes static
	paddingX := int(math.Floor(float64(termSize.X-CANVAS_SIZE.X) / 2))
	paddingY := int(math.Floor(float64(termSize.Y-CANVAS_SIZE.Y) / 2))

	log.Printf("term size: %vx%v", termSize.X, termSize.Y)
	log.Printf("padding: %vx%v", paddingX, paddingY)
	log.Printf("canvas: %vx%v", CANVAS_SIZE.X, CANVAS_SIZE.Y)

	return GameConfig{
		FPS:        FPS,
		DELTA_TIME: time.Duration(int(math.Round(float64(1000)/float64(FPS))) * int(time.Millisecond)),

		CANVAS_SIZE:      CANVAS_SIZE,
		BORDER_THICKNESS: BORDER_THICKNESS,
		TERM_SIZE:        vector2.New(termSize.X, termSize.Y),
		PADDING:          vector2.New(paddingX, paddingY),

		TITLE:  "cURLy.gjt.io",
		FOOTER: "Move: WASD; Restart: R",
	}
}
