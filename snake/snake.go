package snake

import "github.com/edudobay/go-snake/core"

type snakeCell struct {
	pos  int
	next *snakeCell
	prev *snakeCell
}

type snakeCells []int

type Snake struct {
	Cells snakeCells
}

func NewSnake() core.Entity {
	snake := core.NewEntity(EntitySnake)
	snake.AttachComponent(&Snake{})
	snake.AttachComponent(NewPosition())
	return snake
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
	moveResult := s.CheckCollision(direction, board)
	if moveResult != MoveOk {
		return moveResult
	}

	s.GrowHead(direction, board, position)
	s.ShrinkTail(position)

	board.updated()

	return moveResult
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
