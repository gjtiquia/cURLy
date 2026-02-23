package game

import (
	"log"
	"math"
	"time"

	"github.com/gjtiquia/cURLy/internal/vector2"
)

type Config struct {
	FPS       int
	DeltaTime time.Duration

	TermSize        vector2.Type
	BorderThickness vector2.Type
	CanvasSize      vector2.Type
	Padding         vector2.Type

	Title  string
	Footer string
}

func CreateConfig(termSize vector2.Type) Config {
	FPS := 8

	CANVAS_SIZE := vector2.New(20, 8)
	BORDER_THICKNESS := vector2.New(1, 1)

	// uncomment to debug win case
	// FPS = 5
	// CANVAS_SIZE = vector2.New(4, 2)

	// TODO : should make effort to support resize...? right now assumes static
	paddingX := int(math.Floor(float64(termSize.X-CANVAS_SIZE.X) / 2))
	paddingY := int(math.Floor(float64(termSize.Y-CANVAS_SIZE.Y) / 2))

	return Config{
		FPS:       FPS,
		DeltaTime: time.Duration(int(math.Round(float64(1000)/float64(FPS))) * int(time.Millisecond)),

		CanvasSize:      CANVAS_SIZE,
		BorderThickness: BORDER_THICKNESS,
		TermSize:        vector2.New(termSize.X, termSize.Y),
		Padding:         vector2.New(paddingX, paddingY),

		Title:  "cURLy.gjt.io",
		Footer: "Move: WASD; Restart: R",
	}
}
