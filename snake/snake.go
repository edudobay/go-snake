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

func SnakeComponent(entity core.Entity) *Snake {
	return entity.GetComponent(ComponentSnake).(*Snake)
}
