package main

import (
	"strings"
)

func (g *Game) parseCommand(input string) {
	/*input = strings.ToLower(input)*/
	words := strings.Fields(input)

	if len(words) == 0 {
		return
	}

	command := words[0]

	switch command {
	case "help":
		g.output = append(g.output, "commands: help, look, examine, show, clear")
	case "look":
		g.showObjects(&g.currentRoom)
	case "clear":
		g.output = []string{}
	case "show":
		g.showPlayer(&g.player)
	case "examine":
		if len(words) < 2 {
			g.output = append(g.output, "Examine what?")
			return
		}

		target := words[1]

		for i := range g.currentRoom.Objects {
			if g.currentRoom.Objects[i].Name == target {
				g.examineObject(&g.currentRoom.Objects[i])
				return
			}
		}
		g.output = append(g.output, "You don't see that here")
	default:
		if input != "" {
			g.output = append(g.output, "unknown command")
		}
	}
}
