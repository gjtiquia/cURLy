package ansi

import "fmt"

// TODO : test for CMD and PowerShell support, eg. forcing ANSI support in Windows

const RESET = "\x1b[0m" // reset formatting (eg. colors / bold)
const HIDE_CURSOR = "\x1b[?25l"
const SHOW_CURSOR = "\x1b[?25h"
const CURSOR_HOME = "\x1b[H" // moves cursor to top-left corner
const CLEAR = "\x1b[2J"

func ClearAndHideCursor() {
	fmt.Print(CLEAR + HIDE_CURSOR + CURSOR_HOME + RESET)
}

func ClearAndShowCursor() {
	fmt.Print(CLEAR + CURSOR_HOME + SHOW_CURSOR)
}

func ClearAndDrawBuffer(buffer string) {
	fmt.Print(CLEAR + CURSOR_HOME + buffer + RESET)
}
