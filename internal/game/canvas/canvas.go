package canvas

import "github.com/gjtiquia/cURLy/internal/vector2"

type Canvas struct {
	Size  vector2.Type
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

func Create(size vector2.Type) Canvas {
	cells := make([]CellType, size.Y*size.X)
	for i := range cells {
		cells[i] = CellTypeNone
	}
	return Canvas{size, cells}
}

func (this *Canvas) SetCell(pos vector2.Type, cellType CellType) {
	this.Cells[pos.Y*this.Size.X+pos.X] = cellType
}
