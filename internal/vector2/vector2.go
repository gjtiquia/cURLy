package vector2

import (
	"fmt"

	"github.com/gjtiquia/cURLy/internal/random"
)

type vector2 struct {
	X int
	Y int
}

// alias for Vector2
type Type = vector2

func New(x, y int) vector2 {
	return vector2{x, y}
}

func (this vector2) String() string {
	return fmt.Sprintf("(%v, %v)", this.X, this.Y)
}

func (this vector2) Add(other vector2) vector2 {
	return vector2{X: this.X + other.X, Y: this.Y + other.Y}
}

func Random(maxExclusive vector2) vector2 {
	x := random.Range(0, maxExclusive.X)
	y := random.Range(0, maxExclusive.Y)
	return vector2{x, y}
}

var Zero = vector2{0, 0}
var Up = vector2{0, 1}
var Down = vector2{0, -1}
var Left = vector2{-1, 0}
var Right = vector2{1, 0}
