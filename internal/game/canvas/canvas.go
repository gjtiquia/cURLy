package canvas

import "github.com/gjtiquia/cURLy/internal/vector2"
import "github.com/gjtiquia/cURLy/internal/game"

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

func (this *Canvas) SetCellByXY(x, y int, cellType CellType) {
	this.Cells[y*this.Size.X+x] = cellType
}

func CreateCanvas(config game.Config) Canvas {
	canvas := CreateEmptyCanvas(config.TermSize)

	// Define regions (all coordinates are in terminal space, 0-indexed from top-left)
	// Border surrounds the canvas, padding surrounds the border
	borderLeft := config.Padding.X - config.BorderThickness.X
	borderRight := config.Padding.X + config.CanvasSize.X // exclusive for canvas, inclusive for border
	borderTop := config.Padding.Y - config.BorderThickness.Y
	borderBottom := config.Padding.Y + config.CanvasSize.Y // exclusive for canvas, inclusive for border

	canvasLeft := config.Padding.X
	canvasRight := config.Padding.X + config.CanvasSize.X // exclusive
	canvasTop := config.Padding.Y
	canvasBottom := config.Padding.Y + config.CanvasSize.Y // exclusive

	for y := 0; y < config.TermSize.Y; y++ {
		for x := 0; x < config.TermSize.X; x++ {
			// Check if in canvas region
			if x >= canvasLeft && x < canvasRight && y >= canvasTop && y < canvasBottom {
				canvas.SetCellByXY(x, y, CellTypeBg)
				continue
			}

			// Check if in border region (horizontal borders - top and bottom)
			if y >= borderTop && y < canvasTop && x >= borderLeft && x < borderRight+config.BorderThickness.X {
				canvas.SetCellByXY(x, y, CellTypeBorderX)
				continue
			}
			if y >= canvasBottom && y < borderBottom+config.BorderThickness.Y && x >= borderLeft && x < borderRight+config.BorderThickness.X {
				canvas.SetCellByXY(x, y, CellTypeBorderX)
				continue
			}

			// Check if in border region (vertical borders - left and right)
			if y >= canvasTop && y < canvasBottom {
				if x >= borderLeft && x < canvasLeft {
					canvas.SetCellByXY(x, y, CellTypeBorderY)
					continue
				}
				if x >= canvasRight && x < canvasRight+config.BorderThickness.X {
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
