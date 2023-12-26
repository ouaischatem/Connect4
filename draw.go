package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"strconv"
)

// Displays graphics on the screen based on the current state of the game.
func (g *game) Draw(screen *ebiten.Image) {

	screen.Fill(globalBackgroundColor)

	switch g.gameState {
	case titleState:
		g.titleDraw(screen)
	case colorSelectState:
		g.colorSelectDraw(screen)
	case displayColorState:
		g.displayColorDraw(screen)
	case playState:
		g.playDraw(screen)
	case resultState:
		g.resultDraw(screen)
	}
}

// Displays graphics for the title screen.
func (g game) titleDraw(screen *ebiten.Image) {
	text.Draw(screen, "Puissance 4 en réseau", largeFont, 90, 150, globalTextColor)
	text.Draw(screen, "Projet de programmation système", smallFont, 105, 190, globalTextColor)
	text.Draw(screen, "Année 2023-2024", smallFont, 210, 230, globalTextColor)

	if g.stateFrame >= globalBlinkDuration/3 {
		text.Draw(screen, "Appuyez sur entrée", smallFont, 210, 500, globalTextColor)
	}

	if g.isReady() {
		text.Draw(screen, "Votre adversaire est prêt.", smallFont, 160, 650, globalSuccessColor)
	} else {
		text.Draw(screen, "Opposant non connecté.", smallFont, 185, 650, globalErrorColor)
	}
}

// Displays graphics for the player color selection screen.
func (g game) colorSelectDraw(screen *ebiten.Image) {
	text.Draw(screen, "Quelle couleur pour vos pions ?", smallFont, 110, 80, globalTextColor)

	line := 0
	col := 0
	for numColor := 0; numColor < globalNumColor; numColor++ {

		xPos := (globalNumTilesX-globalNumColorCol)/2 + col
		yPos := (globalNumTilesY-globalNumColorLine)/2 + line

		if numColor == g.p1Color {
			vector.DrawFilledCircle(screen, float32(globalTileSize/2+xPos*globalTileSize), float32(globalTileSize+globalTileSize/2+yPos*globalTileSize), globalTileSize/2, globalSelectColor, true)
		}

		if numColor == g.p2Color {
			vector.DrawFilledCircle(screen, float32(globalTileSize/2+xPos*globalTileSize), float32(globalTileSize+globalTileSize/2+yPos*globalTileSize), globalTileSize/2, globalTokenOpacityColors[numColor], true)
			vector.DrawFilledCircle(screen, float32(globalTileSize/2+xPos*globalTileSize), float32(globalTileSize+globalTileSize/2+yPos*globalTileSize), globalTileSize/2-globalCircleMargin, globalTokenOpacityColors[numColor], true)
		} else {
			vector.DrawFilledCircle(screen, float32(globalTileSize/2+xPos*globalTileSize), float32(globalTileSize+globalTileSize/2+yPos*globalTileSize), globalTileSize/2-globalCircleMargin, globalTokenColors[numColor], true)
		}

		col++
		if col >= globalNumColorCol {
			col = 0
			line++
		}
	}
}

// Displays graphics for the screen showing the colors of both players.
func (g game) displayColorDraw(screen *ebiten.Image) {
	text.Draw(screen, "Vos choix de couleurs", smallFont, 200, 80, globalTextColor)
	y := 250

	radius := globalTileSize

	vector.DrawFilledCircle(screen, float32(100+radius), float32(y+radius), float32(radius-globalCircleMargin), globalTokenColors[g.p1Color], true)
	text.Draw(screen, "Vous", smallFont, 165, 500, globalTextColor)

	vector.DrawFilledCircle(screen, float32(400+radius), float32(y+radius), float32(radius-globalCircleMargin), globalTokenColors[g.p2Color], true)
	text.Draw(screen, "Adversaire", smallFont, 425, 500, globalTextColor)

	if g.startTimer > 0 {
		remainingSeconds := ((startCountdown * 60) - g.startTimer) / 60
		if remainingSeconds > 0 {
			text.Draw(screen, "Lancement dans "+strconv.Itoa(remainingSeconds)+" secondes", smallFont, 150, 650, globalSuccessColor)
		} else {
			text.Draw(screen, "Lancement imminent", smallFont, 210, 650, globalSuccessColor)
		}
	} else {
		text.Draw(screen, "En attente de votre adversaire.", smallFont, 135, 650, globalErrorColor)
	}
}

// Displays graphics during the game.
func (g game) playDraw(screen *ebiten.Image) {
	g.drawGrid(screen)

	vector.DrawFilledCircle(screen, float32(globalTileSize/2+g.tokenPosition*globalTileSize), float32(globalTileSize/2), globalTileSize/2-globalCircleMargin, globalTokenColors[g.p1Color], true)
}

// Displays graphics on the results screen.
func (g game) resultDraw(screen *ebiten.Image) {
	g.drawGrid(offScreenImage)

	options := &ebiten.DrawImageOptions{}
	options.ColorScale.ScaleAlpha(0.2)
	screen.DrawImage(offScreenImage, options)

	message := "Égalité"
	if g.result == p1wins {
		message = "Gagné !"
	} else if g.result == p2wins {
		message = "Perdu…"
	}
	text.Draw(screen, message, smallFont, 300, 400, globalTextColor)
}

// Displays the Connect 4 grid, including the already played tokens.
func (g game) drawGrid(screen *ebiten.Image) {
	vector.DrawFilledRect(screen, 0, globalTileSize, globalTileSize*globalNumTilesX, globalTileSize*globalNumTilesY, globalGridColor, true)

	for x := 0; x < globalNumTilesX; x++ {
		for y := 0; y < globalNumTilesY; y++ {

			var tileColor color.Color
			switch g.grid[x][y] {
			case p1Token:
				tileColor = globalTokenColors[g.p1Color]
			case p2Token:
				tileColor = globalTokenColors[g.getP2Color()]
			default:
				tileColor = globalBackgroundColor
			}

			vector.DrawFilledCircle(screen, float32(globalTileSize/2+x*globalTileSize), float32(globalTileSize+globalTileSize/2+y*globalTileSize), globalTileSize/2-globalCircleMargin, tileColor, true)
		}
	}
}
