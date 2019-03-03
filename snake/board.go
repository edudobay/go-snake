package snake

import "fmt"

const BufSize = 16

type signal struct{}

type Board struct {
	width, height int
	cells         []BoardCellType
	snakeHead     *snakeCell
	snakeEnd      *snakeCell
	updates       chan signal
}

type BoardCellType int

type snakeCell struct {
	pos  int
	next *snakeCell
	prev *snakeCell
}

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

	return &Board{map_.width, map_.height, cells, nil, nil, make(chan signal, BufSize)}
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
	b.updates <- signal{}
}

func (b *Board) Updates() <-chan signal {
	return b.updates
}

// putSnake places the snake on the board; the head is placed on the (i, j)
// position, and the tail is arranged linearly in the given direction from the
// head.
func (b *Board) PutSnake(i, j, size int, direction Direction) {
	if b.CellTypeAt(i, j) != BoardCellFree {
		panic("tried to place snake in a non-free position")
	}

	head := b.cellAddress(i, j)
	pos := head

	var prevCell *snakeCell = nil

	for i := 0; i < size; i++ {
		b.cells[pos] = BoardCellSnakeBody

		cell := &snakeCell{pos, nil, prevCell}
		if prevCell != nil {
			prevCell.next = cell
		}
		if b.snakeHead == nil {
			fmt.Printf("new snake head at %v\n", cell)
			b.snakeHead = cell
		}

		prevCell = cell
		b.snakeEnd = cell
		pos += b.step(direction)
	}

	b.updated()
}

func (c *snakeCell) append(pos int) *snakeCell {
	cell := &snakeCell{pos, nil, c}
	c.next = cell
	return cell
}

func (c *snakeCell) prepend(pos int) *snakeCell {
	cell := &snakeCell{pos, c, nil}
	c.prev = cell
	return cell
}

func (c *snakeCell) popNext() *snakeCell {
	next := c.next
	if next != nil {
		c.next = nil
		next.prev = nil
	}
	return next
}

func (b *Board) checkPos(pos int) {
	if pos < 0 || pos >= b.width*b.height {
		panic("board position out of bounds")
	}
}

func (b *Board) growSnakeHead(direction Direction) {
	if b.snakeHead == nil {
		panic("no snake head")
	}
	oldHead := b.snakeHead
	newPos := oldHead.pos + b.step(direction)
	b.checkPos(newPos)

	newHead := b.snakeHead.prepend(newPos)
	b.snakeHead = newHead
	b.cells[newHead.pos] = BoardCellSnakeBody
}

func (b *Board) shrinkSnakeTail() {
	if b.snakeEnd.prev == nil {
		panic("tried to remove only cell (prev == nil)")
	}
	oldEnd := b.snakeEnd
	newEnd := b.snakeEnd.prev
	newEnd.popNext()
	b.snakeEnd = newEnd
	b.cells[oldEnd.pos] = BoardCellFree
}

func (b *Board) GrowSnake(direction Direction) {
	b.growSnakeHead(direction)
	b.updated()
}

func (b *Board) MoveSnake(direction Direction) MoveResult {
	newPos := b.snakeHead.pos + b.step(direction)

	switch b.cells[newPos] {
	case BoardCellSnakeBody:
		if newPos == b.snakeHead.next.pos {
			return MoveSelf
		} else {
			return MoveSelfCollide
		}
	case BoardCellWall:
		return MoveWall
	}

	b.growSnakeHead(direction)
	b.shrinkSnakeTail()
	b.updated()
	return MoveOk
}
