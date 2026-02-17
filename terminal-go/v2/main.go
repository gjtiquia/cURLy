package main

import (
	"log"
	"time"

	"github.com/gdamore/tcell/v3"
	"github.com/gjtiquia/cURLy/terminal-go/v2/internals/logfile"
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
	gameConfig, gameState, canvas, inputBuffer := createGame(s)

	for {
		startTime := time.Now()

		inputBuffer, isExit := DrainTCellEvents(s, inputBuffer)
		if isExit {
			return
		}

		runGameLoop(gameConfig, gameState, canvas, inputBuffer, s)

		elapsedTime := time.Since(startTime)
		remainingTime := gameConfig.DELTA_TIME - elapsedTime
		if remainingTime > 0 {
			time.Sleep(remainingTime)
		}
	}
}

func createGame(s tcell.Screen) (GameConfig, *GameState, GameCanvas, []InputAction) {
	gameConfig := createGameConfig(s.Size())
	gameState := createGameState()
	canvas := createCanvas(gameConfig)
	inputBuffer := make([]InputAction, 0)
	return gameConfig, gameState, canvas, inputBuffer
}

func runGameLoop(gameConfig GameConfig, gameState *GameState, canvas GameCanvas, inputBuffer []InputAction, s tcell.Screen) {
	// game logic
	gameState.onUpdate(gameConfig, inputBuffer)

	// draw
	canvas.resetCanvas(gameConfig)
	gameState.onDraw(gameConfig, canvas)

	// render
	for y, row := range canvas {
		for x, element := range row {
			s.PutStr(x, y, element)
		}
	}
	s.Show()
}
