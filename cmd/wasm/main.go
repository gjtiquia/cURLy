package main

import (
	"fmt"
	"time"

	"github.com/gjtiquia/cURLy/internal/vector2"
)

var canvas Canvas

func main() {
	defer notify(MainExit)

	// termSize := getTermSize()
	termSize := vector2.New(4, 4) // TODO : temp until we figure out js/go bridge
	fmt.Println("termSize", termSize)

	canvas = CreateCanvas(termSize)

	// testing notify
	for i := 1; i <= 10; i++ {
		// canvas.SetCell(vector2.New(i, 0), 0)
		notify(CanvasUpdated)

		time.Sleep(1 * time.Second)
	}
}

type NotifyEvent byte

const (
	MainExit NotifyEvent = iota
	CanvasUpdated
)

type Canvas struct {
	size  vector2.Type
	cells []byte
}

func CreateCanvas(size vector2.Type) Canvas {
	return Canvas{
		size:  size,
		cells: make([]byte, size.Y*size.X),
	}
}

func (this *Canvas) SetCell(pos vector2.Type, value byte) {
	this.cells[pos.Y*this.size.X+pos.X] = value
}
