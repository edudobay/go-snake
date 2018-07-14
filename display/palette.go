package display

type RGBColor struct {
	R, G, B uint8
}

const NumPalettes = 5
var CurrentPalette = 0
var PsychedelicMode = 0

var bgColor = []uint8{180, 191, 155}
var fgColor = []uint8{74, 83, 53}
