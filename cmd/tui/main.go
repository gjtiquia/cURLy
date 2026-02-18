package main

import (
	"log"
	"time"

	"github.com/gdamore/tcell/v3"
	"github.com/gjtiquia/cURLy/internal/logfile"
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
	gameState := CreateGameState(gameConfig.CANVAS_SIZE)

	canvas := createCanvas(gameConfig)
	canvas.drawTitle(gameConfig.TITLE, gameConfig)
	canvas.drawFooter(gameConfig.FOOTER, gameConfig)

	// arbitrary capacity of 4, players probably wont mash more than 4 keys between frames, if so the underlying array should adjust itself
	inputBuffer := make([]InputAction, 0, 4)

	return gameConfig, gameState, canvas, inputBuffer
}

func runGameLoop(gameConfig GameConfig, gameState *GameState, canvas GameCanvas, inputBuffer []InputAction, s tcell.Screen) {
	// game logic
	gameState.OnUpdate(gameConfig, inputBuffer)

	// draw
	canvas.resetCanvas(gameConfig)
	gameState.OnDraw(gameConfig, canvas)

	// render
	for y, row := range canvas {
		for x, element := range row {
			s.PutStr(x, y, element)
		}
	}
	s.Show()
}
