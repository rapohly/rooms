package main

import (
	"strings"
)

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
