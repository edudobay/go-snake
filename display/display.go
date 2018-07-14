package display

import (
	"github.com/edudobay/go-snake/core"
	"github.com/veandco/go-sdl2/sdl"
	"log"
)

const ScreenWidth = 640
const ScreenHeight = 480
const ScreenBpp = 32
const GridWidth = 40
const GridHeight = 30
const BgColor = 0
const FgColor = 1
const NumPalettes = 5

var CurrentPalette = 0
var PsychedelicMode = 0

type Display struct {
	renderer *sdl.Renderer
	sprites *Sprites
}

func (d *Display) Update() {
	d.renderer.Present()
}

func (d *Display) DrawSprite(id Sprite, x, y int32) {
	d.sprites.DrawSprite(id, x, y)
}

func windowResource(window *sdl.Window) core.Resource {
	return core.DisposableResource(func() {
		log.Println("Disposing SDL.Window")
		window.Destroy()
	})
}

func rendererResource(renderer *sdl.Renderer) core.Resource {
	return core.DisposableResource(func() {
		log.Println("Disposing SDL.Renderer")
		renderer.Destroy()
	})
}

func InitDisplay(resources core.HoldsDisposables) (*Display, error) {
	window, err := sdl.CreateWindow(
		"Snake",
		sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		ScreenWidth, ScreenHeight,
		sdl.WINDOW_SHOWN)
	if err != nil {
		return nil, err
	}

	resources.AddDisposable(windowResource(window))

	renderer, err := sdl.CreateRenderer(window, -1, 0)
	if err != nil {
		return nil, err
	}

	resources.AddDisposable(rendererResource(renderer))

	renderer.SetDrawColor(0, 0, 0, 255)
	renderer.Clear()
	renderer.Present()

	// Hack
	sdl.Delay(20)
	renderer.Present()

	sprites, err := LoadSprites(renderer, resources)
	if err != nil {
		panic(err)
	}

	return &Display{renderer, sprites}, nil
}
