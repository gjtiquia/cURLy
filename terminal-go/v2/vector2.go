package main

import "fmt"

type Vector2 struct {
	x int
	y int
}

func (this Vector2) Add(other Vector2) Vector2 {
	return Vector2{x: this.x + other.x, y: this.y + other.y}
}

func (this Vector2) String() string {
	return fmt.Sprintf("(%v, %v)", this.x, this.y)
}
