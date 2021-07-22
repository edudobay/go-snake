package snake

import "github.com/edudobay/go-snake/core"

const BufSize = 16

type signal struct{}

type Board struct {
	width, height int
	system        core.System
	cells         []BoardCellType
	updates       chan signal
}

type BoardCellType int

const (
	BoardCellInvalid BoardCellType = iota
	BoardCellFree
	BoardCellWall
	BoardCellSnakeBody
)

func CreateBoard(map_ GameMap) *Board {
	size := map_.Size()
	cells := make([]BoardCellType, size)
	for i, mapCell := range map_.Cells() {
		cells[i] = cellTypeFromMapCell(mapCell)
	}

	system := core.NewSystem()

	return &Board{map_.width, map_.height, system, cells, make(chan signal, BufSize)}
}

func (b Board) Width() int {
	return b.width
}

func (b Board) Height() int {
	return b.height
}

func (b Board) cellAddress(i, j int) int {
	return i*b.width + j
}

func (b Board) step(direction Direction) int {
	switch direction {
	case Left:
		return -1
	case Right:
		return 1
	case Up:
		return -b.width
	case Down:
		return b.width
	default:
		panic("invalid direction")
	}
}

func (b Board) towards(direction Direction, count int) int {
	return b.step(direction) * count
}

func (b Board) CellTypeAtCellAddress(cell int) BoardCellType {
	if cell < 0 || cell >= len(b.cells) {
		panic("out of board bounds")
	}

	return b.cells[cell]
}

func (b Board) CellTypeAt(i, j int) BoardCellType {
	if i < 0 || i >= b.height || j < 0 || j >= b.width {
		panic("out of board bounds")
	}

	return b.cells[i*b.width+j]
}

func cellTypeFromMapCell(mapCellType MapCellType) BoardCellType {
	switch mapCellType {
	case MapCellInvalid:
		return BoardCellInvalid
	case MapCellFree:
		return BoardCellFree
	case MapCellWall:
		return BoardCellWall
	default:
		panic("invalid map cell type found")
	}
}

func (b Board) Center() (int, int) {
	i := b.height / 2
	j := b.width / 2
	return i, j
}

func (b *Board) updated() {
	for _, component := range b.system.FindComponentsOfType(ComponentPosition) {
		component.(*Position).UpdateBoard(b)
	}
	b.updates <- signal{}
}

func (b *Board) Updates() <-chan signal {
	return b.updates
}

// PutSnake places the snake on the board; the head is placed on the (i, j)
// position, and the tail is arranged linearly in the given direction from the
// head.
func (b *Board) PutSnake(i, j, size int, direction Direction) {
	headAddress := b.cellAddress(i, j)
	if b.CellTypeAtCellAddress(headAddress) != BoardCellFree {
		panic("tried to place snake in a non-free position")
	}

	b.system.AddEntity(NewSnake(headAddress, size, direction, b))
	b.updated()
}

func (b *Board) checkPos(pos int) {
	if pos < 0 || pos >= b.width*b.height {
		panic("board position out of bounds")
	}
}

func (b *Board) UpdateCell(cell int, cellType BoardCellType) {
	b.cells[cell] = cellType
}

func (b *Board) MoveSnake(direction Direction) MoveResult {
	snakeEntity := b.snakeEntity()
	snake := SnakeComponent(snakeEntity)
	position := PositionComponent(snakeEntity)

	return snake.MoveSnake(direction, b, position)
}

func (b *Board) snakeEntity() core.Entity {
	return b.system.FindEntityOrNilById(EntitySnake).(core.Entity)
}
