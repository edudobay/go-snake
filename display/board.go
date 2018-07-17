package display

import (
	"github.com/edudobay/go-snake/snake"
	"github.com/veandco/go-sdl2/sdl"
)

const BoardX = 5
const BoardY = 0

func spriteForBoardCell(cellType snake.BoardCellType) Sprite {
	switch cellType {
	case snake.BoardCellFree:
		return SpriteNone
	case snake.BoardCellWall:
		return SpriteWall
	case snake.BoardCellInvalid:
		return SpriteNone
	default:
		panic("invalid cell type")
	}
}

func (d *Display) DrawBoard(board snake.Board) {
	rect := &sdl.Rect{
		X: BoardX,
		Y: BoardY,
		W: int32(board.Width() * SpriteWidth),
		H: int32(board.Height() * SpriteHeight)}

	SetDrawColorRGB(d.renderer, d.palette.BgColor)
	d.renderer.FillRect(rect)

	for i := 0; i < board.Height(); i++ {
		for j := 0; j < board.Width(); j++ {
			sprite := spriteForBoardCell(board.CellTypeAt(i, j))
			x := BoardX + j * SpriteWidth
			y := BoardY + i * SpriteHeight
			d.DrawSprite(sprite, x, y)
		}
	}
}
