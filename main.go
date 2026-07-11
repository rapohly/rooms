package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	game := &Game{}
	game.Init()

	ebiten.SetWindowSize(screenW, screenH)
	ebiten.SetWindowTitle("ROOMS")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
