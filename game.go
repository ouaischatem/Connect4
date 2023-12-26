package main

import (
	"net"
)

// Data structure to represent the current state of the game.
type game struct {
	gameState     int
	stateFrame    int
	grid          [globalNumTilesX][globalNumTilesY]int
	p1Color       int
	p2Color       int
	turn          int
	tokenPosition int
	result        int

	readyToStart bool
	canPlay      bool
	serverConn   net.Conn
	startTimer   int

	startChan   chan bool
	chosenColor chan bool
	p2Position  chan int

	p2ColorChan chan int

	debug bool
}

// Constants to represent the current game sequence (title screen,
// color selection screen, game, results screen).
const (
	titleState int = iota
	colorSelectState
	displayColorState
	playState
	resultState
)

// Constants to represent tokens in the Connect 4 grid
// (no token, player 1's token, player 2's token).
const (
	noToken int = iota
	p1Token
	p2Token
)

// Constants to represent the turn of the game (player 1 or player 2).
const (
	p1Turn int = iota
	p2Turn
)

// Constants to represent the result of a game
// (draw if the grid is filled without a player winning,
// player 1 winner, or player 2 winner).
const (
	equality int = iota
	p1wins
	p2wins
)

// Resets the game to start a new round.
// The player who lost the last game starts.
func (g *game) reset() {
	for x := 0; x < globalNumTilesX; x++ {
		for y := 0; y < globalNumTilesY; y++ {
			g.grid[x][y] = noToken
		}
	}
}

// Checks if both players have pressed enter and are ready to choose their colors.
// If there is a value in the channel, use it; otherwise, return the previous one.
func (g *game) isReady() bool {
	if g.readyToStart {
		return true
	}

	select {
	case <-g.startChan:
		g.readyToStart = true
		return true
	default:
		return false
	}
}

// Checks if both players have chosen their colors.
// If there is a value in the channel, use it; otherwise, return the previous one.
func (g *game) hasChosenColor() bool {
	if g.canPlay {
		return true
	}

	select {
	case <-g.chosenColor:
		g.canPlay = true
		return true
	default:
		return false
	}
}

// Retrieves the color of Player 2.
// If there is a value in the channel, use it; otherwise, return the previous one.
func (g *game) getP2Color() int {
	color := g.p2Color

	select {
	case color = <-g.p2ColorChan:
		g.p2Color = color
	default:
		// No new value, do not make any changes.
	}

	return g.p2Color
}
