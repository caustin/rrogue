package main

import (
	"github.com/caustin/rrogue/components"
	level2 "github.com/caustin/rrogue/level"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// AutoMoveState tracks the state of auto-movement for smooth progression
type AutoMoveState struct {
	Active        bool
	Direction     struct{ dx, dy int }
	LastMoveTime  time.Time
	MoveCooldown  time.Duration
	StopRequested bool
}

func TakePlayerAction(g *Game) bool {
	players := g.WorldTags["players"]
	turnTaken := false

	// Initialize auto-move state if needed
	if g.AutoMoveState == nil {
		g.AutoMoveState = &AutoMoveState{
			MoveCooldown: 120 * time.Millisecond, // ~8 moves per second
		}
	}

	// Handle ongoing auto-movement
	if g.AutoMoveState.Active {
		return processAutoMovement(g)
	}

	x := 0
	y := 0

	// Check for modifier key (period) + direction for auto-movement
	isAutoMove := ebiten.IsKeyPressed(ebiten.KeyPeriod)

	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		y = -1
		if isAutoMove {
			return startAutoMovement(g, 0, -1)
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		y = 1
		if isAutoMove {
			return startAutoMovement(g, 0, 1)
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		x = -1
		if isAutoMove {
			return startAutoMovement(g, -1, 0)
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		x = 1
		if isAutoMove {
			return startAutoMovement(g, 1, 0)
		}
	}

	// Any key press stops auto-movement
	if g.AutoMoveState.Active && (inpututil.IsKeyJustPressed(ebiten.KeyUp) ||
		inpututil.IsKeyJustPressed(ebiten.KeyDown) ||
		inpututil.IsKeyJustPressed(ebiten.KeyLeft) ||
		inpututil.IsKeyJustPressed(ebiten.KeyRight) ||
		inpututil.IsKeyJustPressed(ebiten.KeyQ) ||
		inpututil.IsKeyJustPressed(ebiten.KeyEscape)) {
		g.AutoMoveState.Active = false
		return false
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		turnTaken = true
	}

	level := g.Map.CurrentLevel

	for _, result := range g.World.Query(players) {
		pos := result.Components[g.Components.Position].(*components.Position)
		index := level.GetIndexFromXY(pos.X+x, pos.Y+y)

		tile := level.Tiles[index]
		if tile.Blocked != true {
			level.Tiles[level.GetIndexFromXY(pos.X, pos.Y)].Blocked = false
			pos.X += x
			pos.Y += y
			level.Tiles[index].Blocked = true
			level.PlayerVisible.Compute(level, pos.X, pos.Y, 8)

		} else if x != 0 || y != 0 {
			if level.Tiles[index].TileType != level2.WALL {
				//Its a tile with a monster -- Fight it
				monsterPosition := components.Position{X: pos.X + x, Y: pos.Y + y}

				AttackSystem(g, pos, &monsterPosition)
			}
		}

	}

	if x != 0 || y != 0 || turnTaken {
		return true
	}
	return false
}

// startAutoMovement initiates auto-movement in the specified direction
func startAutoMovement(g *Game, dx, dy int) bool {
	g.AutoMoveState.Active = true
	g.AutoMoveState.Direction.dx = dx
	g.AutoMoveState.Direction.dy = dy
	g.AutoMoveState.LastMoveTime = time.Now()
	g.AutoMoveState.StopRequested = false

	// Immediately make the first move
	return executePlayerMove(g, dx, dy)
}

// processAutoMovement handles the timed auto-movement logic
func processAutoMovement(g *Game) bool {
	// Check if enough time has passed for the next move
	now := time.Now()
	if now.Sub(g.AutoMoveState.LastMoveTime) < g.AutoMoveState.MoveCooldown {
		return false // Don't take a turn yet
	}

	level := g.Map.CurrentLevel
	players := g.WorldTags["players"]

	for _, result := range g.World.Query(players) {
		pos := result.Components[g.Components.Position].(*components.Position)

		dx := g.AutoMoveState.Direction.dx
		dy := g.AutoMoveState.Direction.dy

		// Check if next position is valid
		nextX := pos.X + dx
		nextY := pos.Y + dy
		nextIndex := level.GetIndexFromXY(nextX, nextY)

		// Stop if we can't move (blocked by wall)
		if level.Tiles[nextIndex].TileType == level2.WALL {
			g.AutoMoveState.Active = false
			return false
		}

		// Stop if there's a monster - attack it
		if level.Tiles[nextIndex].Blocked {
			g.AutoMoveState.Active = false
			monsterPosition := components.Position{X: nextX, Y: nextY}
			AttackSystem(g, pos, &monsterPosition)
			return true
		}

		// Check stop conditions before moving
		if isMonsterVisible(g, level, pos) {
			g.AutoMoveState.Active = false
			return false
		}

		if isAtJunctionOrRoom(level, pos) {
			g.AutoMoveState.Active = false
			return false
		}

		// Execute the move
		g.AutoMoveState.LastMoveTime = now
		return executePlayerMove(g, dx, dy)
	}

	g.AutoMoveState.Active = false
	return false
}

// executePlayerMove performs a single player move in the specified direction
func executePlayerMove(g *Game, dx, dy int) bool {
	players := g.WorldTags["players"]
	level := g.Map.CurrentLevel

	for _, result := range g.World.Query(players) {
		pos := result.Components[g.Components.Position].(*components.Position)
		index := level.GetIndexFromXY(pos.X+dx, pos.Y+dy)

		tile := level.Tiles[index]
		if !tile.Blocked {
			// Move player
			level.Tiles[level.GetIndexFromXY(pos.X, pos.Y)].Blocked = false
			pos.X += dx
			pos.Y += dy
			level.Tiles[index].Blocked = true
			level.PlayerVisible.Compute(level, pos.X, pos.Y, 8)
			return true
		} else if tile.TileType != level2.WALL {
			// Attack monster
			monsterPosition := components.Position{X: pos.X + dx, Y: pos.Y + dy}
			AttackSystem(g, pos, &monsterPosition)
			return true
		}
	}

	return false
}

// isMonsterVisible checks if any monster is visible from the current position
func isMonsterVisible(g *Game, level level2.Level, playerPos *components.Position) bool {
	monsters := g.WorldTags["monsters"]

	for _, monster := range g.World.Query(monsters) {
		monsterPos := monster.Components[g.Components.Position].(*components.Position)
		if level.PlayerVisible.IsVisible(monsterPos.X, monsterPos.Y) {
			return true
		}
	}
	return false
}

// isAtJunctionOrRoom checks if the player is at a corridor junction or in a room
// A junction/room is defined as having more than 2 walkable adjacent tiles
func isAtJunctionOrRoom(level level2.Level, pos *components.Position) bool {
	walkableCount := 0

	// Check all 4 cardinal directions
	directions := []struct{ dx, dy int }{
		{0, -1}, // up
		{0, 1},  // down
		{-1, 0}, // left
		{1, 0},  // right
	}

	for _, dir := range directions {
		checkX := pos.X + dir.dx
		checkY := pos.Y + dir.dy

		// Check bounds
		if !level.InBounds(checkX, checkY) {
			continue
		}

		index := level.GetIndexFromXY(checkX, checkY)
		tile := level.Tiles[index]

		// Count walkable tiles (floor tiles that aren't blocked by monsters)
		if tile.TileType == level2.FLOOR && !tile.Blocked {
			walkableCount++
		}
	}

	// If more than 2 directions are walkable, we're at a junction or in a room
	return walkableCount > 2
}
