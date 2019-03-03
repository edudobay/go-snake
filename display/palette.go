package display

import (
	"github.com/veandco/go-sdl2/sdl"
)

type RGBColor struct {
	R, G, B uint8
}

type Palette struct {
	BgColor RGBColor
	FgColor RGBColor
}

const NumPalettes = 5

var CurrentPalette = 0
var PsychedelicMode = 0

var AllPalettes = []Palette{
	{
		RGBColor{180, 191, 155},
		RGBColor{74, 83, 53},
	},
}

func SetDrawColorRGB(r *sdl.Renderer, color RGBColor) {
	r.SetDrawColor(color.R, color.G, color.B, 255)
}
