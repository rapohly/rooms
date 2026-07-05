package main

import (
	"image"
	"image/color"
	_ "image/png"
	"log"
	"os"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenW        = 1280
	screenH        = 1000
	terminalHeight = 200
)

type Game struct {
	room       *ebiten.Image
	input      string
	output     []string
	lastBlink  time.Time
	showCursor bool
}

func loadImage(path string) (*ebiten.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}

	return ebiten.NewImageFromImage(img), nil
}

func (g *Game) Init() {
	room, err := loadImage("/home/rpohly/projects/rooms/ascii/office.png")
	if err != nil {
		panic(err)
	}

	g.room = room

	g.output = []string{
		"System booting...",
		"Type 'help' for commands.",
	}

	g.lastBlink = time.Now()
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
	if ebiten.IsKeyPressed(ebiten.KeyEnter) {
		g.handleCommand(g.input)
		g.input = ""
	}

	return nil
}

func (g *Game) handleCommand(cmd string) {
	cmd = strings.TrimSpace(strings.ToLower(cmd))

	g.output = append(g.output, "> "+cmd)

	switch cmd {
	case "help":
		g.output = append(g.output, "commands: help, look, clear")
	case "look":
		g.output = append(g.output, "You are in a server room. Fans humming.")
	case "clear":
		g.output = []string{}
	default:
		if cmd != "" {
			g.output = append(g.output, "unknown command")
		}
	}

	if len(g.output) > 8 {
		g.output = g.output[len(g.output)-8:]
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)

	// ----------------------------
	// TOP AREA: ROOM IMAGE
	// ----------------------------

	if g.room != nil {
		op := &ebiten.DrawImageOptions{}

		// scale image to fit width
		imgW, imgH := g.room.Bounds().Dx(), g.room.Bounds().Dy()

		scaleX := float64(screenW) / float64(imgW)
		scaleY := float64(screenH-terminalHeight) / float64(imgH)

		// keep aspect ratio (choose smaller scale)
		scale := scaleX
		if scaleY < scaleX {
			scale = scaleY
		}

		op.GeoM.Scale(scale, scale)

		// center horizontally
		offsetX := (screenW - int(float64(imgW)*scale)) / 2
		offsetY := (screenH - terminalHeight - int(float64(imgH)*scale)) / 2

		op.GeoM.Translate(float64(offsetX), float64(offsetY))

		screen.DrawImage(g.room, op)
	}

	// ----------------------------
	// BOTTOM AREA: TERMINAL
	// ----------------------------

	// terminal background
	ebitenutil.DrawRect(
		screen,
		0,
		screenH-terminalHeight,
		screenW,
		terminalHeight,
		color.RGBA{0, 0, 0, 220},
	)

	// output log
	y := screenH - terminalHeight + 10
	for _, line := range g.output {
		ebitenutil.DebugPrintAt(screen, line, 10, y)
		y += 16
	}

	// input line
	cursor := ""
	if g.showCursor {
		cursor = "_"
	}

	ebitenutil.DebugPrintAt(
		screen,
		"> "+g.input+cursor,
		10,
		screenH-25,
	)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenW, screenH
}

func main() {
	game := &Game{}
	game.Init()

	ebiten.SetWindowSize(screenW, screenH)
	ebiten.SetWindowTitle("ROOMS")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
