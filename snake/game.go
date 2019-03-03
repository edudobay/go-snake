package snake

import (
	"github.com/veandco/go-sdl2/sdl"
)

const BUF_SIZE = 100

type Game struct {
	level     int
	direction Direction
	size      int
	frames    int
	foodCount int
	hasFood   bool
	points    int
	moveCount int

	keyPress  chan sdl.Keysym
}

func NewGame(level int) Game {
	return Game{
		level:     level,
		direction: Down,
		size:      3,
		frames:    0,
		foodCount: 0,
		hasFood:   false,
		points:    0,
		moveCount: 0,
		keyPress:  make(chan sdl.Keysym, BUF_SIZE) }
}

func (g Game) OnKeyPressed(key sdl.Keysym) {
	g.keyPress <- key
}

func (g Game) KeyPresses() chan sdl.Keysym {
	return g.keyPress
}
