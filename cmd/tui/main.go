package main

import (
	"log"
	"time"

	"github.com/gjtiquia/cURLy/internal/game"
	"github.com/gjtiquia/cURLy/internal/logfile"
	"github.com/gjtiquia/cURLy/internal/vector2"
	"github.com/pkg/errors"
)

func main() {
	// logging setup
	err, logPanicAndCloseFile := logfile.Init("log.txt")
	if err != nil {
		log.Panicf("%+v", errors.WithStack(err))
	}
	defer logPanicAndCloseFile()

	// tcell setup
	s, err, finalizeScreen := InitTCellScreen()
	if err != nil {
		log.Panicf("%+v", errors.WithStack(err))
	}
	defer finalizeScreen()

	// gameSetup
	termSize := vector2.New(s.Size())
	gameConfig, gameState, canvas, inputBuffer := game.Create(termSize)

	for {
		startTime := time.Now()

		inputBuffer, isExit := DrainTCellEvents(s, inputBuffer)
		if isExit {
			return
		}

		game.RunLoop(gameConfig, gameState, canvas, inputBuffer)
		RenderCanvas(s, termSize, canvas)

		elapsedTime := time.Since(startTime)
		remainingTime := gameConfig.DeltaTime - elapsedTime
		if remainingTime > 0 {
			time.Sleep(remainingTime)
		}
	}
}
