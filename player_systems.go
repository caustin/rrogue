package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func TakePlayerAction(g *Game) bool {
	players := g.WorldTags["players"]
	turnTaken := false

	x := 0
	y := 0

	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		y = -1
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		y = 1
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		x = -1
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		x = 1
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		turnTaken = true
	}

	level := g.Map.CurrentLevel

	for _, result := range g.World.Query(players) {
		pos := result.Components[position].(*Position)
		index := level.GetIndexFromXY(pos.X+x, pos.Y+y)

		tile := level.Tiles[index]
		if tile.Blocked != true {
			level.Tiles[level.GetIndexFromXY(pos.X, pos.Y)].Blocked = false
			pos.X += x
			pos.Y += y
			level.Tiles[index].Blocked = true
			level.PlayerVisible.Compute(level, pos.X, pos.Y, 8)

		} else if x != 0 || y != 0 {
			if level.Tiles[index].TileType != WALL {
				//Its a tile with a monster -- Fight it
				monsterPosition := Position{X: pos.X + x, Y: pos.Y + y}

				AttackSystem(g, pos, &monsterPosition)
			}
		}

	}

	if x != 0 || y != 0 || turnTaken {
		return true
	}
	return false
}
