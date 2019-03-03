package snake

type Game struct {
	level     int
	direction Direction
	size      int
	frames    int
	foodCount int
	hasFood   bool
	points    int
	moveCount int
	board     *Board
}

func NewGame(level int, board *Board) *Game {
	return &Game{
		level:     level,
		direction: Down,
		size:      3,
		frames:    0,
		foodCount: 0,
		hasFood:   false,
		points:    0,
		moveCount: 0,
		board:     board}
}
