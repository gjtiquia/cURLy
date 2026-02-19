package main

// These function are exported to JavaScript, so can be called using exports.someFunc() in JavaScript.

// must use the //export directive for TinyGo

//export getCanvas
func getCanvas() [][]string {
	return canvas
}
