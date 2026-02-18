package main

import (
	"github.com/gjtiquia/cURLy/terminal-go/v2/internals/random"
	"github.com/gjtiquia/cURLy/terminal-go/v2/internals/vector2"
)

type GameState struct {
	snakeHeadPos     vector2.Type
	snakeBodyPosList []vector2.Type
	snakeDirection   vector2.Type
	foodPos          vector2.Type
}

func CreateGameState(canvasSize vector2.Type) *GameState {
	// random snake head pos
	snakeHeadPos := vector2.Random(canvasSize)

	// random snake direction
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
		snakeHeadPos:   snakeHeadPos,
		snakeDirection: snakeDirection,
	}

	// depends on the existing snake head pos
	gameState.foodPos = gameState.generateRandomFoodPos(canvasSize)

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
	previousSnakeHeadPos := this.snakeHeadPos
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
	ateFood := this.snakeHeadPos == this.foodPos
	if ateFood {
		this.foodPos = this.generateRandomFoodPos(gameConfig.CANVAS_SIZE)

		// add an arbitrary body pos, it will set a new pos anyways when moving body parts forward
		this.snakeBodyPosList = append(this.snakeBodyPosList, vector2.Zero)
	}

	// move each body part forward (move the last one first!)
	if len(this.snakeBodyPosList) > 0 {
		for i := len(this.snakeBodyPosList) - 1; i >= 0; i-- {
			if i == 0 {
				this.snakeBodyPosList[i] = previousSnakeHeadPos
			} else {
				this.snakeBodyPosList[i] = this.snakeBodyPosList[i-1]
			}
		}
	}
}

func (this *GameState) OnDraw(gameConfig GameConfig, canvas GameCanvas) {
	// note: order matters, affects what overlaps what

	canvas.drawCharAtPos(this.foodPos, gameConfig.FOOD_CHAR, gameConfig)

	for _, pos := range this.snakeBodyPosList {
		canvas.drawCharAtPos(pos, gameConfig.SNAKE_BODY_CHAR, gameConfig)
	}

	canvas.drawCharAtPos(this.snakeHeadPos, gameConfig.SNAKE_HEAD_CHAR, gameConfig)
}

func (this *GameState) generateRandomFoodPos(canvasSize vector2.Type) vector2.Type {
	randomFoodPos := vector2.Random(canvasSize)
	isPositionValid := this.isFoodPosValid(randomFoodPos)

	for !isPositionValid {
		randomFoodPos = vector2.Random(canvasSize)
		isPositionValid = this.isFoodPosValid(randomFoodPos)
	}

	return randomFoodPos
}

func (this *GameState) isFoodPosValid(pos vector2.Type) bool {
	if pos == this.snakeHeadPos {
		return false
	}

	for _, bodyPos := range this.snakeBodyPosList {
		if pos == bodyPos {
			return false
		}
	}

	return true
}
