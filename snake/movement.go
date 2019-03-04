package snake

type Direction int

type MoveResult int

const (
	NoDirection Direction = 0
	Up          Direction = 1
	Right       Direction = 2
	Down        Direction = 3
	Left        Direction = 4
)

const (
	MoveOk MoveResult = iota
	MoveSelf
	MoveSelfCollide
	MoveWall
)

func (game *Game) Move(direction Direction) {
	if direction == NoDirection {
		panic("invalid argument for Game.Move: NoDirection")
	}

	game.direction = direction

	game.board.MoveSnake(direction)
}
