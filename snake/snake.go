package snake

import "github.com/edudobay/go-snake/core"

type snakeCell struct {
	pos  int
	next *snakeCell
	prev *snakeCell
}

type snakeCells []int

type Snake struct {
	direction Direction
	Cells     snakeCells
}

func NewSnake(headAddress, size int, direction Direction, b *Board) core.Entity {
	position := NewPosition()

	pos := headAddress
	cells := make(snakeCells, size)

	for i := size - 1; i >= 0; i-- {
		position.UpdateCell(pos, BoardCellSnakeBody)
		cells[i] = pos
		pos += b.step(direction)
	}

	entity := core.NewEntity(EntitySnake)
	entity.AttachComponent(&Snake{Cells: cells, direction: direction})
	entity.AttachComponent(position)
	return entity
}

func (b *Board) Snake() *Snake {
	return b.system.OneComponentOfType(ComponentSnake).(*Snake)
}

func (s *Snake) Type() string {
	return ComponentSnake
}

func (s *Snake) headPos() int {
	return s.Cells[len(s.Cells)-1]
}

func (s *Snake) posFromHead(count int) int {
	return s.Cells[(len(s.Cells)-1)-count]
}

func (s *Snake) MoveSnake(direction Direction, board *Board, position *Position) MoveResult {
	if direction == NoDirection {
		panic("invalid argument for Game.Move: NoDirection")
	}

	s.direction = direction

	moveResult := s.CheckCollision(direction, board)
	if moveResult != MoveOk {
		return moveResult
	}

	s.GrowHead(direction, board, position)
	s.ShrinkTail(position)

	board.updated()

	return moveResult
}

func (s *Snake) AutoMove(board *Board, position *Position) MoveResult {
	return s.MoveSnake(s.direction, board, position)
}

func (s *Snake) CheckCollision(moveDirection Direction, board *Board) MoveResult {
	newPos := s.headPos() + board.step(moveDirection)

	switch board.CellTypeAtCellAddress(newPos) {
	case BoardCellSnakeBody:
		if newPos == s.posFromHead(1) {
			return MoveSelf
		} else {
			return MoveSelfCollide
		}
	case BoardCellWall:
		return MoveWall
	}

	return MoveOk
}

func (s *Snake) GrowHead(direction Direction, board *Board, position *Position) {
	oldHead := s.headPos()
	newHead := oldHead + board.step(direction)
	board.checkPos(newHead)

	s.Cells = append(s.Cells, newHead)

	position.UpdateCell(newHead, BoardCellSnakeBody)
}

func (s *Snake) ShrinkTail(position *Position) {
	if len(s.Cells) <= 1 {
		panic("tried to remove only cell")
	}

	oldEnd := s.Cells[0]
	position.UpdateCell(oldEnd, BoardCellFree)

	s.Cells = s.Cells[1:]
}

func SnakeComponent(entity core.Entity) *Snake {
	return entity.GetComponent(ComponentSnake).(*Snake)
}
