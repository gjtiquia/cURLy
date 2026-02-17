package main

type GameState struct {
	snakeHeadPos   Vector2
	snakeDirection Vector2
}

func createGameState() *GameState {
	gameState := GameState{
		snakeHeadPos:   Vector2{0, 0},
		snakeDirection: Vector2{1, 0},
	}
	return &gameState
}

func (this *GameState) onUpdate(gameConfig GameConfig, inputBuffer []InputAction) {

	// TODO : for now just get the most recent input action
	if len(inputBuffer) > 0 {
		inputAction := inputBuffer[len(inputBuffer)-1]
		// log.Println("action", inputAction)

		switch {
		case inputAction == Up:
			this.snakeDirection = Vector2{0, 1}
		case inputAction == Down:
			this.snakeDirection = Vector2{0, -1}
		case inputAction == Left:
			this.snakeDirection = Vector2{-1, 0}
		case inputAction == Right:
			this.snakeDirection = Vector2{1, 0}
		}
	}

	// update snake head pos
	this.snakeHeadPos = this.snakeHeadPos.Add(this.snakeDirection)

	// wrap around canvas edge
	this.snakeHeadPos.x = this.snakeHeadPos.x % gameConfig.CANVAS_SIZE.x
	this.snakeHeadPos.y = this.snakeHeadPos.y % gameConfig.CANVAS_SIZE.y
	if this.snakeHeadPos.x < 0 {
		this.snakeHeadPos.x += gameConfig.CANVAS_SIZE.x
	}
	if this.snakeHeadPos.y < 0 {
		this.snakeHeadPos.y += gameConfig.CANVAS_SIZE.y
	}
}

func (this GameState) onDraw(gameConfig GameConfig, canvas GameCanvas) {
	canvas.drawChar(this.snakeHeadPos, gameConfig.SNAKE_CHAR, gameConfig)
}
