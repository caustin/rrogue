package main

import (
	"github.com/caustin/rrogue/components"
	"github.com/caustin/rrogue/level"
	"github.com/caustin/rrogue/utils"
	"github.com/hajimehoshi/ebiten/v2"
)

func ProcessRenderables(g *Game, level level.Level, screen *ebiten.Image) {
	for _, result := range g.World.Query(g.WorldTags["renderables"]) {
		pos := result.Components[g.Components.Position].(*components.Position)
		img := result.Components[g.Components.Renderable].(*components.Renderable).Image

		if level.PlayerVisible.IsVisible(pos.X, pos.Y) {
			index := level.GetIndexFromXY(pos.X, pos.Y)
			tile := level.Tiles[index]
			op := utils.GetDrawOptions()

			op.GeoM.Translate(float64(tile.PixelX), float64(tile.PixelY))
			screen.DrawImage(img, op)
			utils.PutDrawOptions(op)
		}

	}
}
