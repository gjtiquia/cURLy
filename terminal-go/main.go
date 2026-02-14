package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gjtiquia/cURLy/terminal-go/ansi"
	"golang.org/x/term"
)

// game config
const FPS = 10
const DELTA_TIME_MS = 1000 / FPS

func main() {
	defer cleanup() // called at the end if no SIGINT or SIGTERM is received
	go listenToSIGINTAndSIGTERM(cleanup)

	file, err := initLogTxt()
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()

	currentTermFd := int(os.Stdout.Fd())
	termWidth, termHeight, err := term.GetSize(currentTermFd)
	if err != nil {
		return
	}

	log.Printf("size: %vx%v", termWidth, termHeight)

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

func initLogTxt() (*os.File, error) {
	// truncate means delete contents on open, create if doesnt exist, write-only
	const fileFlags = os.O_TRUNC | os.O_CREATE | os.O_WRONLY

	// read = 4, write = 2, execute = 1; 6 = 4 + 2 (read write); 0 = octal; 666 = owner/group/others
	const filePerm = 0666

	file, err := os.OpenFile("log.txt", fileFlags, filePerm)
	if err != nil {
		return nil, err
	}

	log.SetOutput(file)
	return file, nil
}

func clearAndDrawBuffer(buffer string) {
	ansi.ClearAndDrawBuffer(buffer)
}

func cleanup() {
	ansi.Cleanup()
}
