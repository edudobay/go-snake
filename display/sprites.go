package display

import (
	"github.com/edudobay/go-snake/core"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type Sprite int

const (
	SpriteInvalid Sprite = iota - 1 // -1
	SpriteNone
	SpriteWall
	SpriteFood
	SpriteTurn
	SpriteBody
	SpriteHead
	SpriteFood2
	SpriteFood3
	SpriteBorderUpperLeft
	SpriteBorderUpperRight
	SpriteBorderLowerRight
	SpriteBorderLowerLeft
	SpriteMarker
	SpriteCount // must be the last entry
)

const SpriteWidth = 11
const SpriteHeight = 11

type Sprites struct {
	texture  *sdl.Texture
	renderer *sdl.Renderer
}

func LoadSprites(renderer *sdl.Renderer, resources core.HoldsDisposables) (*Sprites, error) {
	surface, err := img.Load("data/sprites.png")
	if err != nil {
		return nil, err
	}

	resources.AddDisposable(core.DisposableResource(func() {
		surface.Free()
	}))

	texture, err := renderer.CreateTextureFromSurface(surface)
	if err != nil {
		return nil, err
	}
	return &Sprites{texture, renderer}, nil
}

func (s *Sprites) sourceRectForSprite(id Sprite) *sdl.Rect {
	xOffset := (int32)((id - (SpriteNone + 1)) * SpriteWidth)
	return &sdl.Rect{X: xOffset, Y: 0, W: SpriteWidth, H: SpriteHeight}
}

func (s *Sprites) DrawSprite(id Sprite, x, y int32) {
	if id < SpriteNone || id >= SpriteCount {
		panic("invalid sprite ID")
	}

	if id == SpriteNone { // no-op
		return
	}

	src := s.sourceRectForSprite(id)
	dest := &sdl.Rect{X: x, Y: y, W: SpriteWidth, H: SpriteHeight}

	s.renderer.Copy(s.texture, src, dest)
}
