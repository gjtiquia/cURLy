package main

import "github.com/gjtiquia/cURLy/internals/vector2"

type GameCanvas [][]string

func createCanvas(config GameConfig) GameCanvas {
	canvas := GameCanvas{}

	// Define regions (all coordinates are in terminal space, 0-indexed from top-left)
	// Border surrounds the canvas, padding surrounds the border
	borderLeft := config.PADDING.X - config.BORDER_THICKNESS.X
	borderRight := config.PADDING.X + config.CANVAS_SIZE.X // exclusive for canvas, inclusive for border
	borderTop := config.PADDING.Y - config.BORDER_THICKNESS.Y
	borderBottom := config.PADDING.Y + config.CANVAS_SIZE.Y // exclusive for canvas, inclusive for border

	canvasLeft := config.PADDING.X
	canvasRight := config.PADDING.X + config.CANVAS_SIZE.X // exclusive
	canvasTop := config.PADDING.Y
	canvasBottom := config.PADDING.Y + config.CANVAS_SIZE.Y // exclusive

	for y := 0; y < config.TERM_SIZE.Y; y++ {
		row := []string{}
		for x := 0; x < config.TERM_SIZE.X; x++ {
			// Check if in canvas region
			if x >= canvasLeft && x < canvasRight && y >= canvasTop && y < canvasBottom {
				row = append(row, config.BG_CHAR)
				continue
			}

			// Check if in border region (horizontal borders - top and bottom)
			if y >= borderTop && y < canvasTop && x >= borderLeft && x < borderRight+config.BORDER_THICKNESS.X {
				row = append(row, config.BORDER_X_CHAR)
				continue
			}
			if y >= canvasBottom && y < borderBottom+config.BORDER_THICKNESS.Y && x >= borderLeft && x < borderRight+config.BORDER_THICKNESS.X {
				row = append(row, config.BORDER_X_CHAR)
				continue
			}

			// Check if in border region (vertical borders - left and right)
			if y >= canvasTop && y < canvasBottom {
				if x >= borderLeft && x < canvasLeft {
					row = append(row, config.BORDER_Y_CHAR)
					continue
				}
				if x >= canvasRight && x < canvasRight+config.BORDER_THICKNESS.X {
					row = append(row, config.BORDER_Y_CHAR)
					continue
				}
			}

			// Everything else is padding
			row = append(row, config.PADDING_CHAR)
		}
		canvas = append(canvas, row)
	}

	return canvas
}

func (this GameCanvas) resetCanvas(config GameConfig) {
	for y := 0; y < config.CANVAS_SIZE.Y; y++ {
		for x := 0; x < config.CANVAS_SIZE.X; x++ {
			this.drawChar(x, y, config.BG_CHAR, config)
		}
	}

	// clear message
	y := config.PADDING.Y + config.CANVAS_SIZE.Y + config.BORDER_THICKNESS.Y
	for x := config.PADDING.X; x < config.TERM_SIZE.X; x++ {
		this[y][x] = " "
	}
}

func (this GameCanvas) drawCharAtPos(pos vector2.Type, char string, config GameConfig) {
	this.drawChar(pos.X, pos.Y, char, config)
}

func (this GameCanvas) drawChar(x, y int, char string, config GameConfig) {
	// canvas is drawn from top to bottom but game coordinates is from bottom to top
	// game y=0 is at bottom of canvas, which maps to terminal row: PADDING.y + CANVAS_SIZE.y - 1
	// game y=CANVAS_SIZE.y-1 is at top, which maps to terminal row: PADDING.y
	termY := config.PADDING.Y + (config.CANVAS_SIZE.Y - 1 - y)
	termX := config.PADDING.X + x
	this[termY][termX] = char
}

func (this GameCanvas) drawTitle(title string, config GameConfig) {
	x := config.PADDING.X - config.BORDER_THICKNESS.X
	y := config.PADDING.Y - config.BORDER_THICKNESS.Y - 1

	for i, char := range title {
		if x+i > config.TERM_SIZE.X {
			break
		}

		this[y][x+i] = string(char)
	}
}

func (this GameCanvas) drawMessage(message string, config GameConfig) {
	x := config.PADDING.X - config.BORDER_THICKNESS.X
	y := config.PADDING.Y + config.CANVAS_SIZE.Y + config.BORDER_THICKNESS.Y

	for i, char := range message {
		if x+i >= config.TERM_SIZE.X {
			break
		}

		this[y][x+i] = string(char)
	}
}

func (this GameCanvas) drawFooter(message string, config GameConfig) {
	x := config.PADDING.X - config.BORDER_THICKNESS.X
	y := config.PADDING.Y + config.CANVAS_SIZE.Y + config.BORDER_THICKNESS.Y + 1

	for i, char := range message {
		if x+i >= config.TERM_SIZE.X {
			break
		}

		this[y][x+i] = string(char)
	}
}
