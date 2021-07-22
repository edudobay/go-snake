package snake

import (
	"github.com/edudobay/go-snake/core"
	"sync"
)

type cellUpdate struct {
	Cell     int
	CellType BoardCellType
}

type Position struct {
	Cells       []int
	cellUpdates []cellUpdate
	mutex       sync.Mutex
}

func PositionComponent(entity core.Entity) *Position {
	return entity.GetComponent(ComponentPosition).(*Position)
}

func NewPosition() *Position {
	return &Position{}
}

func (p *Position) Type() string {
	return ComponentPosition
}

func (p *Position) UpdateCell(cell int, cellType BoardCellType) {
	p.mutex.Lock()
	p.cellUpdates = append(p.cellUpdates, cellUpdate{Cell: cell, CellType: cellType})
	p.mutex.Unlock()
}

func (p *Position) UpdateBoard(board *Board) {
	p.mutex.Lock()
	for _, update := range p.cellUpdates {
		board.UpdateCell(update.Cell, update.CellType)
	}
	p.cellUpdates = make([]cellUpdate, 0)
	p.mutex.Unlock()
}
