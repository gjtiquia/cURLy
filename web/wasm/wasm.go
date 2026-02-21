package main

import (
	"fmt"
	"time"

	"github.com/gjtiquia/cURLy/internal/game/canvas"
	"github.com/gjtiquia/cURLy/internal/vector2"
)

var canvasInstance canvas.Type
var inputCh chan byte

func main() {
	defer Notify(MainExit)

	termSize := GetTermSize()
	fmt.Println("termSize", termSize)

	canvasInstance = canvas.Create(termSize)
	inputCh = make(chan byte, 8)

	inputBuffer := make([]byte, 0, 8)

	for i := 0; i < termSize.X; i++ {
		inputBuffer = inputBuffer[:0]
	drain:
		for {
			select {
			case input := <-inputCh:
				inputBuffer = append(inputBuffer, input)
			default:
				break drain
			}
		}

		canvasInstance.SetCell(vector2.New(i, 0), canvas.CellTypeSnakeBody)
		canvasInstance.SetCell(vector2.New(termSize.X-1-i, termSize.Y-1), canvas.CellTypeSnakeBody)
		Notify(CanvasUpdated)

		time.Sleep(1 * time.Second)
	}
}

type NotifyEvent byte

const (
	MainExit NotifyEvent = iota
	CanvasUpdated
)
