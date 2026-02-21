package main

import "github.com/gjtiquia/cURLy/internal/game/canvas"

// These function are exported to JavaScript, so can be called using exports.someFunc() in JavaScript.

// must use the //export directive for TinyGo

//export getCanvasCellsPtr
func getCanvasCellsPtr() *[]canvas.CellType {
	return &canvasInstance.Cells
}

//export onInputAction
func onInputAction(id byte) {
	if inputCh != nil {
		inputCh <- id
	}
}
