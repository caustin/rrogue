package game

import (
	"github.com/caustin/rrogue/level"
	"github.com/caustin/rrogue/utils"
	"github.com/hajimehoshi/ebiten/v2"
)

func ProcessRenderables(g *Game, level level.Level, screen *ebiten.Image) {
	for _, result := range g.World.QueryRenderables() {
		pos := g.World.GetPosition(result)
		img := g.World.GetRenderable(result).Image

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
