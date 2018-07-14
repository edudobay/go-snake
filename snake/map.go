package snake

type GameMap struct {
	width, height int
	cells []MapCellType
}

type MapCellType int8

const (
	MapCellInvalid MapCellType = iota
	MapCellFree 
	MapCellWall 
)
