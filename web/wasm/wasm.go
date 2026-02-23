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

	// TODO : use game.Create and game.RunLoop

	config := game.CreateConfig(termSize)
	canvasInstance = canvas.CreateCanvas(config.TermSize, config.CanvasSize, config.Padding, config.BorderThickness)

	inputCh = make(chan input.Action, 8)
	inputBuffer := input.CreateBuffer()

	for i := 0; i < termSize.X; i++ {
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

		JS_Notify(CanvasUpdated)

		time.Sleep(1 * time.Second)
	}
}

type NotifyEvent byte

const (
	MainExit NotifyEvent = iota
	CanvasUpdated
)
