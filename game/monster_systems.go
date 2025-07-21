package game

import (
	"github.com/caustin/rrogue/components"
	"github.com/caustin/rrogue/level"
	"github.com/norendren/go-fov/fov"
)

func UpdateMonster(game *Game) {
	l := game.Map.CurrentLevel
	playerPosition := components.Position{}

	for _, plr := range game.World.QueryPlayers() {
		pos := game.World.GetPosition(plr)
		playerPosition.X = pos.X
		playerPosition.Y = pos.Y
	}
	for _, result := range game.World.QueryMonsters() {
		pos := game.World.GetPosition(result)
		//mon := result.Components[monster].(*Monster)

		monsterSees := fov.New()
		monsterSees.Compute(l, pos.X, pos.Y, 8)
		if monsterSees.IsVisible(playerPosition.X, playerPosition.Y) {

			if pos.GetManhattanDistance(&playerPosition) == 1 {
				//The monster is right next to the player.  Just smack him down
				game.Systems.Combat.ProcessAttack(pos, &playerPosition)

			} else {
				astar := level.AStar{}
				path := astar.GetPath(l, pos, &playerPosition)
				if len(path) > 1 {
					nextTile := l.Tiles[l.GetIndexFromXY(path[1].X, path[1].Y)]
					if !nextTile.Blocked {
						l.Tiles[l.GetIndexFromXY(pos.X, pos.Y)].Blocked = false
						pos.X = path[1].X
						pos.Y = path[1].Y
						nextTile.Blocked = true
					}
				}
			}

		}

	}

}
