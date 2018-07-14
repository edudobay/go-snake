package snake

type Game struct {
	level int
	direction Direction
	size int
	frames int
	foodCount int
	hasFood bool
	points int
	moveCount int
}

func NewGame(level int) Game {
	return Game{
		level: level,
		direction: Down,
		size: 3,
		frames: 0,
		foodCount: 0,
		hasFood: false,
		points: 0,
		moveCount: 0}
}
