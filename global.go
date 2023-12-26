package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
)

// Constants defining the general parameters of the program.
const (
	globalWidth         = globalNumTilesX * globalTileSize
	globalHeight        = (globalNumTilesY + 1) * globalTileSize
	globalTileSize      = 100
	globalNumTilesX     = 7
	globalNumTilesY     = 6
	globalCircleMargin  = 5
	globalBlinkDuration = 60
	globalNumColorLine  = 3
	globalNumColorCol   = 3
	globalNumColor      = globalNumColorLine * globalNumColorCol

	startCountdown = 10 // En secondes
)

// Constants for message codes.
const (
	codeOpponentConnected      = "#200"
	codeColorsChosen           = "#201"
	codeOpponentColorChoice    = "#301"
	codeOpponentColorSelection = "#302"
	codePlayerCanceledColor    = "#303"
	codeOpponentPosition       = "#401"
)

// Variables defining the general parameters of the program.
var (
	globalBackgroundColor color.Color = color.NRGBA{R: 176, G: 196, B: 222, A: 255}
	globalGridColor       color.Color = color.NRGBA{R: 119, G: 136, B: 153, A: 255}
	globalTextColor       color.Color = color.NRGBA{R: 25, G: 25, B: 5, A: 255}
	globalErrorColor      color.Color = color.NRGBA{R: 255, G: 0, B: 0, A: 255}
	globalSuccessColor    color.Color = color.NRGBA{R: 60, G: 179, B: 113, A: 255}
	globalSelectColor     color.Color = color.NRGBA{R: 25, G: 25, B: 5, A: 255}
	globalNoPlayer        color.Color = color.NRGBA{R: 0, G: 0, B: 0, A: 255}
	smallFont             font.Face
	largeFont             font.Face
	globalTokenColors     = [globalNumColor]color.Color{
		color.NRGBA{R: 255, G: 239, B: 213, A: 255},
		color.NRGBA{R: 60, G: 179, B: 113, A: 255},
		color.NRGBA{R: 154, G: 205, B: 50, A: 255},
		color.NRGBA{R: 189, G: 183, B: 107, A: 255},
		color.NRGBA{R: 255, G: 127, B: 80, A: 255},
		color.NRGBA{R: 240, G: 128, B: 128, A: 255},
		color.NRGBA{R: 152, G: 251, B: 152, A: 255},
		color.NRGBA{R: 221, G: 160, B: 221, A: 255},
		color.NRGBA{R: 245, G: 255, B: 250, A: 255},
	}
	globalTokenOpacityColors = [globalNumColor]color.Color{
		color.NRGBA{R: 255, G: 239, B: 213, A: 64},
		color.NRGBA{R: 60, G: 179, B: 113, A: 64},
		color.NRGBA{R: 154, G: 205, B: 50, A: 64},
		color.NRGBA{R: 189, G: 183, B: 107, A: 64},
		color.NRGBA{R: 255, G: 127, B: 80, A: 64},
		color.NRGBA{R: 240, G: 128, B: 128, A: 64},
		color.NRGBA{R: 152, G: 251, B: 152, A: 64},
		color.NRGBA{R: 221, G: 160, B: 221, A: 64},
		color.NRGBA{R: 245, G: 255, B: 250, A: 64},
	}

	globalBackgroundColorChoice = [7]color.Color{
		color.NRGBA{R: 176, G: 196, B: 222, A: 255},
		color.NRGBA{R: 255, G: 255, B: 153, A: 255},
		color.NRGBA{R: 153, G: 255, B: 255, A: 255},
		color.NRGBA{R: 153, G: 153, B: 255, A: 255},
		color.NRGBA{R: 255, G: 153, B: 255, A: 255},
		color.NRGBA{R: 255, G: 153, B: 204, A: 255},
		color.NRGBA{R: 224, G: 224, B: 224, A: 255},
	}
	colorIndex     int
	offScreenImage *ebiten.Image
)
