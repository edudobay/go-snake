package snake

type Board struct {
	width, height int
	cells         []BoardCellType
}

type BoardCellType int

const (
	BoardCellInvalid BoardCellType = iota
	BoardCellFree
	BoardCellWall
)

func CreateBoard(map_ GameMap) Board {
	size := map_.Size()
	cells := make([]BoardCellType, size)
	for i, mapCell := range map_.Cells() {
		cells[i] = cellTypeFromMapCell(mapCell)
	}

	return Board{map_.width, map_.height, cells}
}

func (b Board) Width() int {
	return b.width
}

func (b Board) Height() int {
	return b.height
}

func (b Board) CellTypeAt(i, j int) BoardCellType {
	if i < 0 || i >= b.height || j < 0 || j >= b.width {
		panic("out of board bounds")
	}

	return b.cells[i*b.width+j]
}

func cellTypeFromMapCell(mapCellType MapCellType) BoardCellType {
	switch mapCellType {
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
