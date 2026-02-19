package main

// These function are exported to JavaScript, so can be called using exports.someFunc() in JavaScript.

func multiply(x, y int) int {
	return x * y
}
