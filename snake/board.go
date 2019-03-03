package snake

type Board struct {
	map_ GameMap
}

type BoardCellType int

const (
	BoardCellInvalid BoardCellType = iota
	BoardCellFree
	BoardCellWall
)

func CreateBoard(map_ GameMap) Board {
	return Board{map_}
}

func (b Board) Width() int {
	return b.map_.width
}

func (b Board) Height() int {
	return b.map_.height
}

func (b Board) CellTypeAt(i, j int) BoardCellType {
	switch b.map_.CellTypeAt(i, j) {
	case MapCellInvalid:
		return BoardCellInvalid
	case MapCellFree:
		return BoardCellFree
	case MapCellWall:
		return BoardCellWall
	default:
		panic("invalid map cell type found")
	}
}
