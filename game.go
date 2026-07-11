package main

import (
	_ "image/png"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct {
	currentRoom Room
	input       string
	output      []string
	lastBlink   time.Time
	showCursor  bool
	player      Player
}

func (g *Game) Init() {
	img, err := loadImage("assets/office.png")
	if err != nil {
		panic(err)
	}

	g.output = []string{
		"System booting...",
		"Type 'help' for commands.",
	}

	g.lastBlink = time.Now()

	g.currentRoom = Room{
		Name:  "Office",
		Image: img,
		Objects: []Object{
			{Name: "Lamp", HasKey: true},
			{Name: "Bookshelf", HasKey: false},
			{Name: "File Cabinet", HasKey: false},
		},
		Exits: []Exit{
			{Name: "Door", IsUnlocked: false, NextRoom: nil},
		},
	}
}

func (g *Game) Update() error {
	// blink cursor
	if time.Since(g.lastBlink) > 500*time.Millisecond {
		g.showCursor = !g.showCursor
		g.lastBlink = time.Now()
	}

	// text input
	for _, r := range ebiten.AppendInputChars(nil) {
		g.input += string(r)
	}

	// backspace
	if ebiten.IsKeyPressed(ebiten.KeyBackspace) && len(g.input) > 0 {
		g.input = g.input[:len(g.input)-1]
	}

	// enter
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		g.handleCommand(g.input)
		g.input = ""
	}

	return nil
}
