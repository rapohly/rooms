package main

func (g *Game) showObjects(room *Room) {
	for _, obj := range room.Objects {
		g.output = append(g.output, obj.Name)
	}
}

func (g *Game) showPlayer(player *Player) {
	for _, item := range player.Inventory {
		g.output = append(g.output, item.Name)
	}
}

func (g *Game) examineObject(obj *Object) {
	for _, item := range obj.Items {
		if item.Name == "key" {
			g.output = append(g.output, "THIS OBJECT HAS THE KEY!")
			g.player.Inventory = append(g.player.Inventory, item)
			return
		}
	}
	g.output = append(g.output, "NO KEY HERE")
}
