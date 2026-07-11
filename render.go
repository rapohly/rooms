package main

import (
	"bytes"
	"embed"
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenW        = 1280
	screenH        = 1000
	terminalHeight = 200
)

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
