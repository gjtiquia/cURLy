package main

import (
	"fmt"
	"time"
)

// game config
const FPS = 10
const DELTA_TIME_MS = 1000 / FPS

func main() {

	// game loop
	for i := 0; i < 10; i++ {
		clearAndDrawBuffer("hello world\n")
		time.Sleep(DELTA_TIME_MS * time.Millisecond)
	}

	cleanup()
}

// ANSI
// TODO : test for CMD and PowerShell support, eg. forcing ANSI support in Windows
const RESET = "\x1b[0m"
const HIDE_CURSOR = "\x1b[?25l"
const SHOW_CURSOR = "\x1b[?25h"
const CURSOR_HOME = "\x1b[H"
const CLEAR = "\x1b[2J"

func clearAndDrawBuffer(buffer string) {
	fmt.Print(CLEAR + HIDE_CURSOR + CURSOR_HOME + buffer + RESET)
}

func cleanup() {
	const clear = true

	if clear {
		fmt.Print(CLEAR + CURSOR_HOME + SHOW_CURSOR)
	} else {
		fmt.Print(SHOW_CURSOR)
	}
}
