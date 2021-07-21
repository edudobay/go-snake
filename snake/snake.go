package snake

import "github.com/edudobay/go-snake/core"

const ComponentSnake = "snake"

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
	snake := core.NewEntity("snake")
	snake.AttachComponent(&Snake{})
	return snake
}

func (b *Board) Snake() *Snake {
	return b.system.OneComponentOfType(ComponentSnake).(*Snake)
}

func (s *Snake) Type() string {
	return ComponentSnake
}
