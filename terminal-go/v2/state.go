package main

import "github.com/gjtiquia/cURLy/terminal-go/v2/internals/vector2"

type GameState struct {
	snakeHeadPos   vector2.Type
	snakeDirection vector2.Type
}

func createGameState() *GameState {
	gameState := GameState{
		snakeHeadPos:   vector2.New(0, 0),
		snakeDirection: vector2.New(1, 0),
	}
	return &gameState
}

func (this *GameState) onUpdate(gameConfig GameConfig, inputBuffer []InputAction) {

	// TODO : for now just get the most recent input action
	if len(inputBuffer) > 0 {
		inputAction := inputBuffer[len(inputBuffer)-1]
		// log.Println("action", inputAction)

		switch {
		case inputAction == Up && this.snakeDirection != vector2.Down:
			this.snakeDirection = vector2.Up
		case inputAction == Down && this.snakeDirection != vector2.Up:
			this.snakeDirection = vector2.Down
		case inputAction == Left && this.snakeDirection != vector2.Right:
			this.snakeDirection = vector2.Left
		case inputAction == Right && this.snakeDirection != vector2.Left:
			this.snakeDirection = vector2.Right
		}
	}

	// update snake head pos
	this.snakeHeadPos = this.snakeHeadPos.Add(this.snakeDirection)

	// wrap around canvas edge
	this.snakeHeadPos.X = this.snakeHeadPos.X % gameConfig.CANVAS_SIZE.X
	this.snakeHeadPos.Y = this.snakeHeadPos.Y % gameConfig.CANVAS_SIZE.Y
	if this.snakeHeadPos.X < 0 {
		this.snakeHeadPos.X += gameConfig.CANVAS_SIZE.X
	}
	if this.snakeHeadPos.Y < 0 {
		this.snakeHeadPos.Y += gameConfig.CANVAS_SIZE.Y
	}
}

func (this GameState) onDraw(gameConfig GameConfig, canvas GameCanvas) {
	canvas.drawCharAtPos(this.snakeHeadPos, gameConfig.SNAKE_CHAR, gameConfig)
}
