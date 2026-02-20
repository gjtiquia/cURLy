package main

import "github.com/gjtiquia/cURLy/internal/vector2"

// These functions are imported from JavaScript, as they doesn't define a body.
// You should define these functions in the WebAssembly 'env' module from JavaScript.

// must use the //export directive for TinyGo

//export getTermSize
func getTermSize() vector2.Type

//export notify
func notify(event NotifyEvent)
