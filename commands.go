package main

func (g *Game) showObjects(room *Room) {
	for _, obj := range room.Objects {
		g.output = append(g.output, obj.Name)
	}
}
