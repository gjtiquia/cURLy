package main

import (
	"fmt"
	"time"

	"github.com/gjtiquia/cURLy/internal/vector2"
)

var canvas Canvas

func main() {
	defer Notify(MainExit)

	termSize := GetTermSize()
	fmt.Println("termSize", termSize)

	canvas = CreateCanvas(termSize, ' ')

	// testing notify
	for i := 0; i < termSize.X; i++ {
		canvas.SetCell(vector2.New(i, 1), byte('x'))
		Notify(CanvasUpdated)

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

func CreateCanvas(size vector2.Type, defaultCell byte) Canvas {
	cells := make([]byte, size.Y*size.X)
	for i := range cells {
		cells[i] = defaultCell
	}
	return Canvas{size, cells}
}

func (this *Canvas) SetCell(pos vector2.Type, value byte) {
	this.cells[pos.Y*this.size.X+pos.X] = value
}
