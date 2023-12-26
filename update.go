package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Updating the game state based on keyboard inputs.
func (g *game) Update() error {
	g.stateFrame++
	g.backgroundColorUpdate()

	switch g.gameState {
	case titleState:
		if g.titleUpdate() {
			g.gameState++
		}
	case colorSelectState:
		if g.colorSelectUpdate() {
			g.communicatToServer(codeOpponentColorChoice)
			g.gameState++
		}
	case displayColorState:
		if g.displayColorUpdate() {
			g.gameState++
		}
	case playState:
		g.tokenPosUpdate()
		var lastXPositionPlayed = -1
		var lastYPositionPlayed = -1

		var p2Position int

		select {
		case p2Position = <-g.p2Position:
			lastXPositionPlayed, lastYPositionPlayed = g.p2Update(p2Position)
		default:
			if g.turn == p1Turn {
				lastXPositionPlayed, lastYPositionPlayed = g.p1Update()
			}
		}

		if lastXPositionPlayed >= 0 {
			finished, result := g.checkGameEnd(lastXPositionPlayed, lastYPositionPlayed)
			if finished {
				g.result = result
				g.gameState++
			}
		}
	case resultState:
		if g.resultUpdate() {
			g.reset()
			g.gameState = playState
		}
	}

	return nil
}

// Updating the game state on the color display screen.
func (g *game) displayColorUpdate() bool {
	if !g.hasChosenColor() {
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			g.gameState--
			g.startTimer = 0
			g.communicatToServer(codePlayerCanceledColor)
			return false
		}
		g.getP2Color()
		return false
	}

	g.startTimer++
	return g.startTimer >= (startCountdown * 60)
}

// Updating the background color when the "C" key is pressed.
func (g *game) backgroundColorUpdate() {
	if inpututil.IsKeyJustPressed(ebiten.KeyC) {
		colorIndex = (colorIndex + 1) % len(globalBackgroundColorChoice)
		globalBackgroundColor = globalBackgroundColorChoice[colorIndex]
	}
}

// Updating the game state on the title screen.
func (g *game) titleUpdate() bool {
	g.stateFrame = g.stateFrame % globalBlinkDuration
	return g.isReady() && inpututil.IsKeyJustPressed(ebiten.KeyEnter)
}

// Updating the game state during color selection.
func (g *game) colorSelectUpdate() bool {

	col := g.p1Color % globalNumColorCol
	line := g.p1Color / globalNumColorLine

	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		col = (col + 1) % globalNumColorCol
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		col = (col - 1 + globalNumColorCol) % globalNumColorCol
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		line = (line + 1) % globalNumColorLine
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		line = (line - 1 + globalNumColorLine) % globalNumColorLine
	}

	newColor := line*globalNumColorLine + col

	if g.getP2Color() == newColor {
		g.p1Color = newColor
		return false
	}

	if g.p1Color != newColor {
		g.p1Color = newColor
		g.communicatToServer(codeOpponentColorSelection)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		if g.p1Color != g.getP2Color() {
			return true
		}
	}

	return false
}

// Handling the position of the next token to be played by Player 1.
func (g *game) tokenPosUpdate() {
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		g.tokenPosition = (g.tokenPosition - 1 + globalNumTilesX) % globalNumTilesX
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		g.tokenPosition = (g.tokenPosition + 1) % globalNumTilesX
	}
}

// Handling the moment when the next token is played by Player 1.
func (g *game) p1Update() (int, int) {
	lastXPositionPlayed := -1
	lastYPositionPlayed := -1
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		if updated, yPos := g.updateGrid(p1Token, g.tokenPosition); updated {
			g.turn = p2Turn
			lastXPositionPlayed = g.tokenPosition
			lastYPositionPlayed = yPos
			g.communicatToServer(codeOpponentPosition)
		}
	}
	return lastXPositionPlayed, lastYPositionPlayed
}

// Handling the position of the next token played by Player 2 and when this token is played.
func (g *game) p2Update(position int) (int, int) {
	updated, yPos := g.updateGrid(p2Token, position)
	for ; !updated; updated, yPos = g.updateGrid(p2Token, position) {
		position = (position + 1) % globalNumTilesX
	}
	g.turn = p1Turn
	return position, yPos
}

// Updating the game state on the results screen.
func (g game) resultUpdate() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyEnter)
}

// Updating the game grid when a token is inserted into
// the column with coordinate (x) position.
func (g *game) updateGrid(token, position int) (updated bool, yPos int) {
	for y := globalNumTilesY - 1; y >= 0; y-- {
		if g.grid[position][y] == noToken {
			updated = true
			yPos = y
			g.grid[position][y] = token
			return
		}
	}
	return
}

// Checking for the end of the game: Did the last player who placed a token win?
// Is the grid filled without a winner (draw)? Or should the game continue?
func (g game) checkGameEnd(xPos, yPos int) (finished bool, result int) {

	tokenType := g.grid[xPos][yPos]

	// horizontal
	count := 0
	for x := xPos; x < globalNumTilesX && g.grid[x][yPos] == tokenType; x++ {
		count++
	}
	for x := xPos - 1; x >= 0 && g.grid[x][yPos] == tokenType; x-- {
		count++
	}

	if count >= 4 {
		if tokenType == p1Token {
			return true, p1wins
		}
		return true, p2wins
	}

	// vertical
	count = 0
	for y := yPos; y < globalNumTilesY && g.grid[xPos][y] == tokenType; y++ {
		count++
	}

	if count >= 4 {
		if tokenType == p1Token {
			return true, p1wins
		}
		return true, p2wins
	}

	// diag haut gauche/bas droit
	count = 0
	for x, y := xPos, yPos; x < globalNumTilesX && y < globalNumTilesY && g.grid[x][y] == tokenType; x, y = x+1, y+1 {
		count++
	}

	for x, y := xPos-1, yPos-1; x >= 0 && y >= 0 && g.grid[x][y] == tokenType; x, y = x-1, y-1 {
		count++
	}

	if count >= 4 {
		if tokenType == p1Token {
			return true, p1wins
		}
		return true, p2wins
	}

	// diag haut droit/bas gauche
	count = 0
	for x, y := xPos, yPos; x >= 0 && y < globalNumTilesY && g.grid[x][y] == tokenType; x, y = x-1, y+1 {
		count++
	}

	for x, y := xPos+1, yPos-1; x < globalNumTilesX && y >= 0 && g.grid[x][y] == tokenType; x, y = x+1, y-1 {
		count++
	}

	if count >= 4 {
		if tokenType == p1Token {
			return true, p1wins
		}
		return true, p2wins
	}

	// equality ?
	if yPos == 0 {
		for x := 0; x < globalNumTilesX; x++ {
			if g.grid[x][0] == noToken {
				return
			}
		}
		return true, equality
	}

	return
}
