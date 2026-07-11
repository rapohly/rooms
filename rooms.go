package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Item struct {
	Name string
}

type Object struct {
	Name  string
	Items []Item
}

type Exit struct {
	Name       string
	IsUnlocked bool
	NextRoom   *Room
}

type Room struct {
	Name    string
	Image   *ebiten.Image
	Objects []Object
	Exits   []Exit
}
