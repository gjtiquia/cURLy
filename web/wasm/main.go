package main

import (
	"time"

	"github.com/gjtiquia/cURLy/internal/game"
	"github.com/gjtiquia/cURLy/internal/game/canvas"
	"github.com/gjtiquia/cURLy/internal/game/input"
)

var canvasInstance canvas.Type
var inputCh chan input.Action

func main() {
	defer JS_Notify(MainExit)

	termSize := JS_GetTermSize()
	config, gameState, c, inputBuffer := game.Create(termSize)

	// globals
	canvasInstance = c
	inputCh = make(chan input.Action, 8)

	// game loop
	for {
		startTime := time.Now()

		inputBuffer = inputBuffer[:0]
	drain:
		for {
			select {
			case inputAction := <-inputCh:
				inputBuffer = append(inputBuffer, inputAction)
			default:
				break drain
			}
		}

		game.RunLoop(config, gameState, c, inputBuffer)
		JS_Notify(CanvasUpdated)

		elapsedTime := time.Since(startTime)
		remainingTime := config.DeltaTime - elapsedTime
		if remainingTime > 0 {
			time.Sleep(remainingTime)
		}
	}
}

type NotifyEvent byte

const (
	MainExit NotifyEvent = iota
	CanvasUpdated
)
