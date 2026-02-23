package game

import (
	"github.com/gjtiquia/cURLy/internal/game/canvas"
	"github.com/gjtiquia/cURLy/internal/game/input"
	"github.com/gjtiquia/cURLy/internal/vector2"
)

func CreateGame(termSize vector2.Type) (Config, *GameState, canvas.Type, input.Buffer) {
	config := CreateConfig(termSize)
	state := CreateGameState(config.CanvasSize)

	c := canvas.CreateCanvas(config.TermSize, config.CanvasSize, config.Padding, config.BorderThickness)
	c.DrawTitle(config.Title, config.Padding, config.BorderThickness, config.TermSize)
	c.DrawFooter(config.Footer, config.Padding, config.BorderThickness, config.TermSize, config.CanvasSize)

	// arbitrary capacity of 4, players probably wont mash more than 4 keys between frames, if so the underlying array should adjust itself
	inputBuffer := input.CreateBuffer()

	return config, state, c, inputBuffer
}

func RunGameLoop(config Config, gameState *GameState, c canvas.Type, inputBuffer input.Buffer) {
	// game logic
	gameState.OnUpdate(config, inputBuffer)

	// draw
	c.ResetCanvas(config.TermSize, config.CanvasSize, config.Padding, config.BorderThickness)
	gameState.OnDraw(config, c)
}
