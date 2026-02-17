package main

type GameCanvas [][]string

func createCanvas(config GameConfig) GameCanvas {
	canvas := GameCanvas{}

	// Define regions (all coordinates are in terminal space, 0-indexed from top-left)
	// Border surrounds the canvas, padding surrounds the border
	borderLeft := config.PADDING.x - config.BORDER_THICKNESS.x
	borderRight := config.PADDING.x + config.CANVAS_SIZE.x // exclusive for canvas, inclusive for border
	borderTop := config.PADDING.y - config.BORDER_THICKNESS.y
	borderBottom := config.PADDING.y + config.CANVAS_SIZE.y // exclusive for canvas, inclusive for border

	canvasLeft := config.PADDING.x
	canvasRight := config.PADDING.x + config.CANVAS_SIZE.x // exclusive
	canvasTop := config.PADDING.y
	canvasBottom := config.PADDING.y + config.CANVAS_SIZE.y // exclusive

	for y := 0; y < config.TERM_SIZE.y; y++ {
		row := []string{}
		for x := 0; x < config.TERM_SIZE.x; x++ {
			// Check if in canvas region
			if x >= canvasLeft && x < canvasRight && y >= canvasTop && y < canvasBottom {
				row = append(row, config.BG_CHAR)
				continue
			}

			// Check if in border region (horizontal borders - top and bottom)
			if y >= borderTop && y < canvasTop && x >= borderLeft && x < borderRight+config.BORDER_THICKNESS.x {
				row = append(row, config.BORDER_X_CHAR)
				continue
			}
			if y >= canvasBottom && y < borderBottom+config.BORDER_THICKNESS.y && x >= borderLeft && x < borderRight+config.BORDER_THICKNESS.x {
				row = append(row, config.BORDER_X_CHAR)
				continue
			}

			// Check if in border region (vertical borders - left and right)
			if y >= canvasTop && y < canvasBottom {
				if x >= borderLeft && x < canvasLeft {
					row = append(row, config.BORDER_Y_CHAR)
					continue
				}
				if x >= canvasRight && x < canvasRight+config.BORDER_THICKNESS.x {
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
	for y := 0; y < config.CANVAS_SIZE.y; y++ {
		for x := 0; x < config.CANVAS_SIZE.x; x++ {
			this.drawChar(Vector2{x, y}, config.BG_CHAR, config)
		}
	}
}

func (this GameCanvas) drawChar(position Vector2, char string, config GameConfig) {
	// canvas is drawn from top to bottom but game coordinates is from bottom to top
	// game y=0 is at bottom of canvas, which maps to terminal row: PADDING.y + CANVAS_SIZE.y - 1
	// game y=CANVAS_SIZE.y-1 is at top, which maps to terminal row: PADDING.y
	termY := config.PADDING.y + (config.CANVAS_SIZE.y - 1 - position.y)
	termX := config.PADDING.x + position.x
	this[termY][termX] = char
}
