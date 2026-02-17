package main

import (
	"github.com/gjtiquia/cURLy/terminal-go/v2/internals/random"
	"github.com/gjtiquia/cURLy/terminal-go/v2/internals/vector2"
)

type GameState struct {
	snakeHeadPos   vector2.Type
	snakeDirection vector2.Type
	foodPos        vector2.Type
}

func CreateGameState(canvasSize vector2.Type) *GameState {

	snakeDirection := vector2.Zero
	switch random.Range(0, 4) {
	case 0:
		snakeDirection = vector2.Up
	case 1:
		snakeDirection = vector2.Down
	case 2:
		snakeDirection = vector2.Left
	case 3:
		snakeDirection = vector2.Right
	}

	gameState := GameState{
		snakeHeadPos:   vector2.Random(canvasSize),
		snakeDirection: snakeDirection,
	}

	gameState.foodPos = gameState.generateRandomFoodPos(canvasSize)
	// gameState.foodPos = vector2.New(6, 3)

	return &gameState
}

func (this *GameState) OnUpdate(gameConfig GameConfig, inputBuffer []InputAction) {
	// update snake direction
	if len(inputBuffer) > 0 {
		// TODO : for now just get the most recent input action
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

	// ate food handling and spawn new food handling
	if this.snakeHeadPos == this.foodPos {
		this.foodPos = this.generateRandomFoodPos(gameConfig.CANVAS_SIZE)
	}
}

func (this *GameState) OnDraw(gameConfig GameConfig, canvas GameCanvas) {
	canvas.drawCharAtPos(this.foodPos, gameConfig.FOOD_CHAR, gameConfig)
	canvas.drawCharAtPos(this.snakeHeadPos, gameConfig.SNAKE_CHAR, gameConfig)
}

func (this *GameState) generateRandomFoodPos(canvasSize vector2.Type) vector2.Type {

	// TODO : account for body

	randomFoodPos := this.snakeHeadPos
	for randomFoodPos == this.snakeHeadPos {
		randomFoodPos = vector2.Random(canvasSize)
	}
	return randomFoodPos
}
