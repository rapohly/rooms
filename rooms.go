package main

import (
	"bytes"
	"embed"
	"image"
	"image/color"
	_ "image/png"
	"log"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	screenW        = 1280
	screenH        = 1000
	terminalHeight = 200
)

type Player struct {
	Inventory []Object
	Health    int
}

type Object struct {
	Name   string
	HasKey bool
}

type Exit struct {
	IsUnlocked bool
}

type Room struct {
	Name    string
	Image   *ebiten.Image
	Objects []Object
	Exits   []Exit
}

type Game struct {
	currentRoom Room
	input       string
	output      []string
	lastBlink   time.Time
	showCursor  bool
	player      Player
}

func (g *Game) showObjects(room *Room) {
	for _, obj := range room.Objects {
		g.output = append(g.output, obj.Name)
	}
}

//go:embed assets/**
var assets embed.FS

func loadImage(path string) (*ebiten.Image, error) {
	f, err := assets.ReadFile(path)
	if err != nil {
		return nil, err
	}

	img, _, err := image.Decode(bytes.NewReader(f))
	if err != nil {
		return nil, err
	}

	return ebiten.NewImageFromImage(img), nil
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

func (g *Game) handleCommand(cmd string) {
	cmd = strings.TrimSpace(strings.ToLower(cmd))

	g.output = append(g.output, "> "+cmd)

	switch cmd {
	case "help":
		g.output = append(g.output, "commands: help, look, clear")
	case "look":
		g.output = append(g.output, "You are in an office.")
		g.showObjects(&g.currentRoom)
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

	if g.currentRoom.Image != nil {
		op := &ebiten.DrawImageOptions{}

		// scale image to fit width
		imgW, imgH := g.currentRoom.Image.Bounds().Dx(), g.currentRoom.Image.Bounds().Dy()

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

		screen.DrawImage(g.currentRoom.Image, op)
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
	y := screenH - terminalHeight + 20
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
