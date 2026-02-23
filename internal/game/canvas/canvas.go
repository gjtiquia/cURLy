package canvas

import (
	"fmt"

	"github.com/gjtiquia/cURLy/internal/vector2"
)

type Canvas struct {
	Size vector2.Type

	// use 1D array rather than 2D array of arrays
	// More performant, all in one block, and also easier to parse via pointer and length
	Cells []CellType
}
type Type = Canvas

type CellType byte

const (
	CellTypeNone      CellType = ' '
	CellTypePadding   CellType = ' '
	CellTypeBorderX   CellType = '-'
	CellTypeBorderY   CellType = '|'
	CellTypeBg        CellType = ' '
	CellTypeSnakeHead CellType = '0'
	CellTypeSnakeBody CellType = 'o'
	CellTypeFood      CellType = '*'
)

func CreateEmptyCanvas(size vector2.Type) Canvas {
	cells := make([]CellType, size.Y*size.X)
	for i := range cells {
		cells[i] = CellTypeNone
	}
	return Canvas{size, cells}
}

func (this *Canvas) SetCellByPos(pos vector2.Type, cellType CellType) {
	this.SetCellByXY(pos.X, pos.Y, cellType)
}

func (this *Canvas) SetCellByXY(x, y int, cellType CellType) error {
	index := y*this.Size.X + x
	if index >= len(this.Cells) {
		return fmt.Errorf("Canvas.SetCellByXYRaw: (%d, %d) greater than size %v!", x, y, this.Size)
	}

	this.Cells[y*this.Size.X+x] = cellType
	return nil
}

func (this *Canvas) SetCellByXYRaw(x, y int, rawByte byte) error {
	index := y*this.Size.X + x
	if index >= len(this.Cells) {
		return fmt.Errorf("Canvas.SetCellByXYRaw: (%d, %d) greater than size %v!", x, y, this.Size)
	}

	this.Cells[index] = CellType(rawByte)
	return nil
}

// TODO : kinda ugly, better accept interface instead?

func (this *Canvas) ResetCanvas(termSize, canvasSize, padding, borderThickness vector2.Type) {
	for y := 0; y < canvasSize.Y; y++ {
		for x := 0; x < canvasSize.X; x++ {
			this.SetCellByXY(x, y, CellTypeBg)
		}
	}

	// clear message
	y := padding.Y + canvasSize.Y + borderThickness.Y
	for x := padding.X; x < termSize.X; x++ {
		this.SetCellByXY(x, y, CellTypePadding)
	}
}

func CreateCanvas(termSize, canvasSize, padding, borderThickness vector2.Type) Canvas {
	canvas := CreateEmptyCanvas(termSize)

	// Define regions (all coordinates are in terminal space, 0-indexed from top-left)
	// Border surrounds the canvas, padding surrounds the border
	borderLeft := padding.X - borderThickness.X
	borderRight := padding.X + canvasSize.X // exclusive for canvas, inclusive for border
	borderTop := padding.Y - borderThickness.Y
	borderBottom := padding.Y + canvasSize.Y // exclusive for canvas, inclusive for border

	canvasLeft := padding.X
	canvasRight := padding.X + canvasSize.X // exclusive
	canvasTop := padding.Y
	canvasBottom := padding.Y + canvasSize.Y // exclusive

	for y := 0; y < termSize.Y; y++ {
		for x := 0; x < termSize.X; x++ {
			// Check if in canvas region
			if x >= canvasLeft && x < canvasRight && y >= canvasTop && y < canvasBottom {
				canvas.SetCellByXY(x, y, CellTypeBg)
				continue
			}

			// Check if in border region (horizontal borders - top and bottom)
			if y >= borderTop && y < canvasTop && x >= borderLeft && x < borderRight+borderThickness.X {
				canvas.SetCellByXY(x, y, CellTypeBorderX)
				continue
			}
			if y >= canvasBottom && y < borderBottom+borderThickness.Y && x >= borderLeft && x < borderRight+borderThickness.X {
				canvas.SetCellByXY(x, y, CellTypeBorderX)
				continue
			}

			// Check if in border region (vertical borders - left and right)
			if y >= canvasTop && y < canvasBottom {
				if x >= borderLeft && x < canvasLeft {
					canvas.SetCellByXY(x, y, CellTypeBorderY)
					continue
				}
				if x >= canvasRight && x < canvasRight+borderThickness.X {
					canvas.SetCellByXY(x, y, CellTypeBorderY)
					continue
				}
			}

			// Everything else is padding
			canvas.SetCellByXY(x, y, CellTypePadding)
		}
	}

	return canvas
}

func (this *Canvas) DrawTitle(title string, padding, borderThickness, termSize vector2.Type) {
	x := padding.X - borderThickness.X
	y := padding.Y - borderThickness.Y - 1

	for i, char := range title {
		if x+i > termSize.X {
			break
		}

		this.SetCellByXYRaw(x+i, y, byte(char))
	}
}

func (this *Canvas) DrawMessage(message string, padding, borderThickness, termSize, canvasSize vector2.Type) {
	x := padding.X - borderThickness.X
	y := padding.Y + canvasSize.Y + borderThickness.Y

	for i, char := range message {
		if x+i >= termSize.X {
			break
		}

		this.SetCellByXYRaw(x+i, y, byte(char))
	}
}

func (this *Canvas) DrawFooter(message string, padding, borderThickness, termSize, canvasSize vector2.Type) {
	x := padding.X - borderThickness.X
	y := padding.Y + canvasSize.Y + borderThickness.Y + 1

	for i, char := range message {
		if x+i >= termSize.X {
			break
		}

		this.SetCellByXYRaw(x+i, y, byte(char))
	}
}
