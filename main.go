package main

import (
	"flag"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"golang.org/x/image/font/opentype"
	"log"
	"net"
)

// Setting up the fonts used for display.
func init() {
	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}

	smallFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size: 30,
		DPI:  72,
	})
	if err != nil {
		log.Fatal(err)
	}

	largeFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size: 50,
		DPI:  72,
	})
	if err != nil {
		log.Fatal(err)
	}
}

// Creating an auxiliary image for displaying results.
func init() {
	offScreenImage = ebiten.NewImage(globalWidth, globalHeight)
}

// Creation, configuration, and launch of the game.
func main() {
	ip := flag.String("ip", "", "Adresse IP du serveur")
	port := flag.String("port", "", "Port du serveur")
	flag.Parse()

	if *ip == "" || *port == "" {
		log.Fatal("L'adresse IP/Port est obligatoire. Exemple : -ip 127.0.0.1 -port 8080")
	}

	conn, err := connectToServer(*ip, *port)
	if err != nil {
		log.Println("Erreur de connexion:", err)
		return
	}
	defer conn.Close()

	g := game{}
	g.debug = false // Enables or disables debug mode.
	g.serverConn = conn
	g.startChan = make(chan bool, 1)
	g.chosenColor = make(chan bool, 1)
	g.p2Position = make(chan int, 1)
	g.p2ColorChan = make(chan int, 1)

	go g.handleServer()

	ebiten.SetWindowTitle("Programmation système : projet puissance 4")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(&g); err != nil {
		log.Fatal(err)
	}
}

// Function to establish a connection with the server.
func connectToServer(ip, port string) (net.Conn, error) {
	conn, err := net.Dial("tcp", ip+":"+port)
	if err != nil {
		return nil, err
	}
	log.Println("Client connecté avec succès à", ip+":"+port)
	return conn, nil
}
