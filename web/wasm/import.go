package main

import (
	"github.com/gjtiquia/cURLy/internal/vector2"
)

// These functions are imported from JavaScript, as they doesn't define a body.
// You should define these functions in the WebAssembly 'env' module from JavaScript.

// must use the //export directive for TinyGo

// JS_ prefix is to easily see that we are calling imported JS functions

//export getTermSize
func getTermSize(ptr *vector2.Type)
func JS_GetTermSize() vector2.Type { // a wrapper
	var termSize vector2.Type
	getTermSize(&termSize)
	return termSize
}

//export notify
func notify(event byte)
func JS_Notify(event NotifyEvent) {
	notify(byte(event))
}
