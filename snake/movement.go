package snake

import (
	"log"
)

type Direction int

const (
	NoDirection Direction = 0
	Up          Direction = 1
	Right       Direction = 2
	Down        Direction = 3
	Left        Direction = 4
)

func (game *Game) Move(direction Direction) {
	if direction == NoDirection {
		panic("invalid argument for Game.Move: NoDirection")
	}
	log.Printf("Moving to direction %v\n", direction)
	game.direction = direction

	game.board.MoveSnake(direction)
}
