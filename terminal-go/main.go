package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// game config
const FPS = 10
const DELTA_TIME_MS = 1000 / FPS

func main() {
	defer cleanup() // called at the end if no SIGINT or SIGTERM is received
	go listenToSIGINTAndSIGTERM(cleanup)

	// game loop
	for { // "'while' is spelled 'for' in Go"
		clearAndDrawBuffer("hello world\n")
		time.Sleep(DELTA_TIME_MS * time.Millisecond)
	}
}

func listenToSIGINTAndSIGTERM(cleanupFunc func()) {
	// create a channel, type os.Signal, buffer 1 (required by signal.Notify)
	channel := make(chan os.Signal, 1)

	// notify channel on os.Interrupt (SIGINT) or SIGTERM
	signal.Notify(channel, os.Interrupt, syscall.SIGTERM)

	// blocked until receives a notification from channel
	<-channel
	// receivedSignal := <-channel // go does not allow unused variables
	// fmt.Println("receivedSignal", receivedSignal)

	// need to manually cleanup cuz deferred calls will not be called after exit
	cleanupFunc()

	// caught the signal myself, so also need to exit myself, as it overrides the default behavior
	os.Exit(0)
}

// ANSI
// TODO : test for CMD and PowerShell support, eg. forcing ANSI support in Windows
const RESET = "\x1b[0m"
const HIDE_CURSOR = "\x1b[?25l"
const SHOW_CURSOR = "\x1b[?25h"
const CURSOR_HOME = "\x1b[H"
const CLEAR = "\x1b[2J"

func clearAndDrawBuffer(buffer string) {
	const hideCursor = true
	// const hideCursor = false

	if hideCursor {
		fmt.Print(CLEAR + HIDE_CURSOR + CURSOR_HOME + buffer + RESET)
	} else {
		fmt.Print(CLEAR + CURSOR_HOME + buffer + RESET)
	}
}

func cleanup() {
	// fmt.Println("cleanup")

	const clear = true
	// const clear = false

	if clear {
		fmt.Print(CLEAR + CURSOR_HOME + SHOW_CURSOR)
	} else {
		fmt.Print(SHOW_CURSOR)
	}
}
