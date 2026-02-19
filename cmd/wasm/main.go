package main

import (
	"fmt"
	"time"
)

var canvas [][]string = nil

func main() {
	defer notify("main-exit")

	termSize := getTermSize()
	fmt.Println("termSize", termSize)

	canvas = make([][]string, termSize.Y)

	// testing notify
	for i := 1; i <= 10; i++ {
		canvas = append(canvas, []string{fmt.Sprintf("row %d", i)})
		notify("canvas-updated")

		time.Sleep(1 * time.Second)
	}
}
