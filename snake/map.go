package snake

type GameMap struct {
	width, height int
	cells         []MapCellType
}

type MapCellType int8

const (
	MapCellInvalid MapCellType = iota
	MapCellFree
	MapCellWall
)

func (m GameMap) CellTypeAt(i, j int) MapCellType {
	if i < 0 || i >= m.height || j < 0 || j >= m.width {
		panic("out of map bounds")
	}

	return m.cells[i*m.width+j]
}

func (m GameMap) Cells() []MapCellType {
	return m.cells
}

func (m GameMap) Size() int {
	return m.width * m.height
}
