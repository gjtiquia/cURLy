package ansi

import "fmt"

// TODO : test for CMD and PowerShell support, eg. forcing ANSI support in Windows
const isAnsiSupported = true

const RESET = "\x1b[0m"
const HIDE_CURSOR = "\x1b[?25l"
const SHOW_CURSOR = "\x1b[?25h"
const CURSOR_HOME = "\x1b[H"
const CLEAR = "\x1b[2J"

func IsAnsiSupported() bool {
	return isAnsiSupported
}

func ClearAndDrawBuffer(buffer string) {
	if !isAnsiSupported {
		return
	}

	const hideCursor = true
	// const hideCursor = false

	if hideCursor {
		fmt.Print(CLEAR + HIDE_CURSOR + CURSOR_HOME + buffer + RESET)
	} else {
		fmt.Print(CLEAR + CURSOR_HOME + buffer + RESET)
	}
}

func Cleanup() {
	if !isAnsiSupported {
		return
	}

	// fmt.Println("cleanup")

	const clear = true
	// const clear = false

	if clear {
		fmt.Print(CLEAR + CURSOR_HOME + SHOW_CURSOR)
	} else {
		fmt.Print(SHOW_CURSOR)
	}
}
