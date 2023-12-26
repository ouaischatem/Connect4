package main

import (
	"bufio"
	"log"
	"strconv"
	"strings"
)

// Function to manage interactions with the server
func (g *game) handleServer() {
	reader := bufio.NewReader(g.serverConn)

	for {
		rawMessage, err := reader.ReadString('\n')
		if err != nil {
			continue
		}

		message := strings.ReplaceAll(strings.TrimSpace(rawMessage), " ", "")

		if message == "" {
			continue
		}

		var debugMessage string

		switch {
		case strings.HasPrefix(message, codeOpponentColorChoice):
			if color, err := strconv.Atoi(strings.TrimPrefix(message, codeOpponentColorChoice)); err == nil {
				g.p2ColorChan <- color
			}
		case strings.HasPrefix(message, codeOpponentColorSelection):
			if selectColor, err := strconv.Atoi(strings.TrimPrefix(message, codeOpponentColorSelection)); err == nil {
				g.p2ColorChan <- selectColor
			}
		case strings.HasPrefix(message, codeOpponentPosition):
			if position, err := strconv.Atoi(strings.TrimPrefix(message, codeOpponentPosition)); err == nil {
				g.p2Position <- position
			}
		case message == codeOpponentConnected:
			g.startChan <- true
			debugMessage = "L'adversaire est connecté, il est possible de lancer le jeu."
		case message == codeColorsChosen:
			g.chosenColor <- true
			debugMessage = "Les couleurs ont été choisies des deux côtés."
		}

		if g.debug && debugMessage == "" {
			log.Println("{RECEIVE} <- " + message)
		} else if g.debug {
			log.Println(debugMessage)
		}
	}
}

// Function to send a message to the server
func (g *game) sendServerMessage(message string) {
	writer := bufio.NewWriter(g.serverConn)
	_, err := writer.WriteString(message + "\n")
	if err != nil {
		return
	}
	err = writer.Flush()
	if g.debug {
		log.Println("{SEND} -> " + message)
	}
}

// Communicates the game state to the server based on the provided code
func (g *game) communicatToServer(code string) {
	var message string

	switch code {
	case codeOpponentColorChoice:
		message = code + " " + strconv.Itoa(g.p1Color)
	case codePlayerCanceledColor:
		message = code
	case codeOpponentColorSelection:
		message = code + " " + strconv.Itoa(g.p1Color)
	case codeOpponentPosition:
		message = codeOpponentPosition + " " + strconv.Itoa(g.tokenPosition)
	}

	g.sendServerMessage(message)
}
