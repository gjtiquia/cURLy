package game

import (
	"fmt"
	"slices"

	"github.com/gjtiquia/cURLy/internal/game/canvas"
	"github.com/gjtiquia/cURLy/internal/game/input"
	"github.com/gjtiquia/cURLy/internal/random"
	"github.com/gjtiquia/cURLy/internal/vector2"
	"github.com/pkg/errors"
)

type PlayState int

const (
	GamePlaying PlayState = iota
	GameLost
	GameWon
)

type GameState struct {
	playState PlayState

	snakeHeadPos     vector2.Type
	snakeBodyPosList []vector2.Type
	snakeDirection   vector2.Type

	foodPos vector2.Type
	score   int

	remainingInputBuffer []input.Action
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
		playState: GamePlaying,

		snakeHeadPos:     snakeHeadPos,
		snakeDirection:   snakeDirection,
		snakeBodyPosList: make([]vector2.Type, 0, canvasSize.X*canvasSize.Y),

		foodPos: vector2.Zero,
		score:   0,

		remainingInputBuffer: make([]input.Action, 0, 4),
	}

	// depends on the existing snake head pos
	foodPos, err := gameState.generateRandomFoodPos(canvasSize)
	if err != nil {
		panic(err) // should not happen because canvas is empty
	}

	gameState.foodPos = foodPos
	return &gameState
}

func (this *GameState) OnUpdate(config Config, inputBuffer []input.Action) {

	if slices.Contains(inputBuffer, input.ActionRestart) {
		// Overwrite the struct at our pointer so the caller sees the new state; assigning to `this` would only change our local copy.
		// `this` is a pointer to gameState
		// `this = CreateGameState()` simply changes the local pointer to point to a new struct, but the pointer of the caller still points to the existing struct, so no visual change
		// `*this` dereferences the pointer, so we have access to the struct
		// `*this = *CreateGameState` means we override the struct fields in the original struct, keeping the same pointer address
		*this = *CreateGameState(config.CanvasSize)
		return
	}

	if this.playState != GamePlaying {
		return
	}

	// update snake direction from input
	inputAction := input.ActionNone

	switch {
	case len(inputBuffer) > 0:
		inputAction = inputBuffer[0]

		// resets remaining input buffer and add the rest of input buffer into remaining input buffer
		this.remainingInputBuffer = append(this.remainingInputBuffer[:0], inputBuffer[1:]...)

	case len(this.remainingInputBuffer) > 0:
		inputAction = this.remainingInputBuffer[0]

		this.remainingInputBuffer = this.remainingInputBuffer[1:] // [1:] is safe in Go even if length is 0
	}

	switch {
	case inputAction == input.ActionUp && this.snakeDirection != vector2.Down:
		this.snakeDirection = vector2.Up

	case inputAction == input.ActionDown && this.snakeDirection != vector2.Up:
		this.snakeDirection = vector2.Down

	case inputAction == input.ActionLeft && this.snakeDirection != vector2.Right:
		this.snakeDirection = vector2.Left

	case inputAction == input.ActionRight && this.snakeDirection != vector2.Left:
		this.snakeDirection = vector2.Right
	}

	previousSnakeHeadPos := this.snakeHeadPos
	nextSnakeHeadPos := this.snakeHeadPos.Add(this.snakeDirection)

	if this.isOverlappingBody(nextSnakeHeadPos) {
		this.playState = GameLost
		return
	}

	// update snake head pos
	this.snakeHeadPos = nextSnakeHeadPos

	// wrap around canvas edge
	this.snakeHeadPos.X = this.snakeHeadPos.X % config.CanvasSize.X
	this.snakeHeadPos.Y = this.snakeHeadPos.Y % config.CanvasSize.Y
	if this.snakeHeadPos.X < 0 {
		this.snakeHeadPos.X += config.CanvasSize.X
	}
	if this.snakeHeadPos.Y < 0 {
		this.snakeHeadPos.Y += config.CanvasSize.Y
	}

	// move each body part forward (move the last one first!)
	previousLastBodyPos := this.snakeHeadPos // fallback
	if len(this.snakeBodyPosList) > 0 {
		previousLastBodyPos = this.snakeBodyPosList[len(this.snakeBodyPosList)-1]
		for i := len(this.snakeBodyPosList) - 1; i >= 0; i-- {
			if i == 0 {
				this.snakeBodyPosList[i] = previousSnakeHeadPos
			} else {
				this.snakeBodyPosList[i] = this.snakeBodyPosList[i-1]
			}
		}
	}

	// ate food handling and spawn new food handling
	ateFood := this.snakeHeadPos == this.foodPos
	if ateFood {
		this.score += 10 // add 10 seems happier than add 1 lol
		this.snakeBodyPosList = append(this.snakeBodyPosList, previousLastBodyPos)

		// order matters! generate food AFTER new body positions are updated
		foodPos, err := this.generateRandomFoodPos(config.CanvasSize)
		if err == nil { // will be nil if canvas is full
			this.foodPos = foodPos
		}
	}

	if this.checkWin(config.CanvasSize) {
		this.playState = GameWon
	}
}

func (this *GameState) checkWin(canvasSize vector2.Type) bool {
	return 1+len(this.snakeBodyPosList) == canvasSize.X*canvasSize.Y
}

func (this *GameState) OnDraw(config Config, c canvas.Type) {
	// note: order matters, affects what overlaps what

	c.SetCellByCanvasPos(this.foodPos, canvas.CellTypeFood, config.Padding, config.CanvasSize)

	for _, pos := range this.snakeBodyPosList {
		c.SetCellByCanvasPos(pos, canvas.CellTypeSnakeBody, config.Padding, config.CanvasSize)
	}

	c.SetCellByCanvasPos(this.snakeHeadPos, canvas.CellTypeSnakeHead, config.Padding, config.CanvasSize)

	switch this.playState {
	case GamePlaying:
		c.DrawMessage(fmt.Sprintf("Score: %v", this.score), config.Padding, config.BorderThickness, config.TermSize, config.CanvasSize)
	case GameLost:
		c.DrawMessage(fmt.Sprintf("You Lost! Score: %v", this.score), config.Padding, config.BorderThickness, config.TermSize, config.CanvasSize)
	case GameWon:
		c.DrawMessage(fmt.Sprintf("You Won! Score: %v", this.score), config.Padding, config.BorderThickness, config.TermSize, config.CanvasSize)
	default:
		c.DrawMessage("", config.Padding, config.BorderThickness, config.TermSize, config.CanvasSize)
	}
}

func (this *GameState) generateRandomFoodPos(canvasSize vector2.Type) (vector2.Type, error) {
	if this.checkWin(canvasSize) {
		return vector2.Zero, errors.New("no space to generate food pos")
	}

	randomFoodPos := vector2.Random(canvasSize)
	isPositionValid := this.isFoodPosValid(randomFoodPos)

	for !isPositionValid {
		randomFoodPos = vector2.Random(canvasSize)
		isPositionValid = this.isFoodPosValid(randomFoodPos)
	}

	return randomFoodPos, nil
}

func (this *GameState) isFoodPosValid(pos vector2.Type) bool {
	if pos == this.snakeHeadPos {
		return false
	}

	if this.isOverlappingBody(pos) {
		return false
	}

	return true
}

func (this *GameState) isOverlappingBody(pos vector2.Type) bool {
	return slices.Contains(this.snakeBodyPosList, pos)
}
