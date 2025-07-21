package game

import (
	"github.com/caustin/rrogue/config"
	"github.com/caustin/rrogue/events"
	"github.com/caustin/rrogue/systems"
	"github.com/caustin/rrogue/world"
	"github.com/hajimehoshi/ebiten/v2"
)

// Game holds all data the entire game will need.
type Game struct {
	Map           GameMap
	World         world.WorldService
	EventBus      *events.EventBus
	Systems       *systems.SystemRegistry
	GameData      config.GameData
	Turn          TurnState
	TurnCounter   int
	AutoMoveState *AutoMoveState
}

// NewGame creates a new Game Object and initializes the data
func NewGame() *Game {
	g := &Game{}
	g.Map = NewGameMap()
	g.GameData = config.NewGameData()

	// Create world service
	g.World = world.NewGameWorld(g.Map.CurrentLevel)

	// Create event bus
	g.EventBus = events.NewEventBus()

	// Create and initialize all systems
	g.Systems = systems.NewSystemRegistry(g.World, g.EventBus)

	// Register all event handlers
	g.Systems.RegisterAllHandlers()

	// Wire up GameStateSystem references to Game struct fields
	if g.Systems.GameState != nil {
		// Cast TurnState to int pointer for GameStateSystem
		turnPtr := (*int)(&g.Turn)
		g.Systems.GameState.SetGameReferences(turnPtr, &g.TurnCounter)
	}

	// Wire up MapBridge with tile unblocking function
	if g.Systems.MapBridge != nil {
		unblockTileFunc := func(x, y int) {
			level := g.Map.CurrentLevel
			tile := level.Tiles[level.GetIndexFromXY(x, y)]
			tile.Blocked = false
		}
		g.Systems.MapBridge.SetGameReference(unblockTileFunc)
	}

	// Temporary event handlers are no longer needed -
	// MapBridge handles tile cleanup and GameStateSystem handles game over

	g.Turn = WaitingForPlayerInput
	g.TurnCounter = 0
	return g
}

// Update is called each tic.
func (g *Game) Update() error {
	switch g.Turn {
	case WaitingForPlayerInput:
		if TakePlayerAction(g) {
			// Publish turn change event instead of direct assignment
			if g.Systems.GameState != nil {
				g.Systems.GameState.ChangeTurn(systems.ProcessingMonsterTurn)
				g.Systems.GameState.IncrementTurn()
			} else {
				// Fallback for during migration
				g.Turn = ProcessingMonsterTurn
				g.TurnCounter++
			}
		}
	case ProcessingMonsterTurn:
		UpdateMonster(g)
		// Publish turn change event instead of direct assignment
		if g.Systems.GameState != nil {
			g.Systems.GameState.ChangeTurn(systems.WaitingForPlayerInput)
		} else {
			// Fallback for during migration
			g.Turn = WaitingForPlayerInput
		}
	default:
		panic("unhandled default case")
	}

	return nil

}

// Draw is called each draw cycle and is where we will blit.
func (g *Game) Draw(screen *ebiten.Image) {
	//Draw the Map
	level := g.Map.CurrentLevel
	level.DrawLevel(screen, g.GameData)
	ProcessRenderables(g, level, screen)
	ProcessUserLog(g, screen)
	ProcessHUD(g, screen)

}

// Layout will return the screen dimensions.
func (g *Game) Layout(w, h int) (int, int) {
	return g.GameData.TileWidth * g.GameData.ScreenWidth, g.GameData.TileHeight * g.GameData.ScreenHeight
}
