package main

import (
	"time"

	"github.com/gjtiquia/cURLy/internal/game/canvas"
	"github.com/gjtiquia/cURLy/internal/game/input"
	"github.com/gjtiquia/cURLy/internal/vector2"
)

var canvasInstance canvas.Type
var inputCh chan input.Action

func main() {
	defer JS_Notify(MainExit)

	termSize := JS_GetTermSize()
	canvasInstance = canvas.Create(termSize)

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

		canvasInstance.SetCell(vector2.New(i, 0), canvas.CellTypeSnakeBody)
		canvasInstance.SetCell(vector2.New(termSize.X-1-i, termSize.Y-1), canvas.CellTypeSnakeBody)
		JS_Notify(CanvasUpdated)

		time.Sleep(1 * time.Second)
	}
}

type NotifyEvent byte

const (
	MainExit NotifyEvent = iota
	CanvasUpdated
)
