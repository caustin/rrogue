package systems

import (
	"github.com/caustin/rrogue/events"
	"github.com/caustin/rrogue/world"
)

// SystemRegistry manages all event-driven systems and their lifecycle
type SystemRegistry struct {
	Combat     *CombatSystem
	GameState  *GameStateSystem
	Map        *MapSystem
	GameBridge *GameBridge
	MapBridge  *MapBridge
	UI         *UISystem

	// Dependencies
	world    world.WorldService
	eventBus *events.EventBus
}

// NewSystemRegistry creates and initializes all systems with their dependencies
func NewSystemRegistry(world world.WorldService, eventBus *events.EventBus) *SystemRegistry {
	registry := &SystemRegistry{
		world:    world,
		eventBus: eventBus,
	}

	// Create systems with dependencies
	registry.Combat = NewCombatSystem(world, eventBus)
	registry.GameBridge = NewGameBridge(eventBus)
	registry.MapBridge = NewMapBridge(eventBus)
	registry.UI = NewUISystem(world, eventBus)

	// Create GameStateSystem
	registry.GameState = NewGameStateSystem(world, eventBus)

	// For now, create placeholder systems - we'll wire them up properly later
	// registry.Map = NewMapSystem(eventBus, world, &MapAdapter{})

	return registry
}

// RegisterAllHandlers subscribes all systems to their respective events
func (r *SystemRegistry) RegisterAllHandlers() {
	r.Combat.RegisterHandlers()
	r.UI.RegisterHandlers()

	if r.GameState != nil {
		r.GameState.RegisterHandlers()
	}

	// Register other systems when they're implemented
	// if r.Map != nil {
	//     r.Map.RegisterHandlers()
	// }
}

// GameAdapter implements GameStateAdapter interface for the Game struct
type GameAdapter struct {
	gameRef interface{} // Will hold reference to Game struct
}

// NewGameAdapter creates a new game adapter
func NewGameAdapter() *GameAdapter {
	return &GameAdapter{}
}

// SetGameReference allows the adapter to reference the main game
func (ga *GameAdapter) SetGameReference(game interface{}) {
	ga.gameRef = game
}

// SetTurnState sets the turn state in the game
func (ga *GameAdapter) SetTurnState(state interface{}) {
	if ga.gameRef != nil {
		// Type assertion to access Game struct
		if gamePtr, ok := ga.gameRef.(*interface{}); ok {
			if game, ok := (*gamePtr).(interface{ SetTurnState(interface{}) }); ok {
				game.SetTurnState(state)
			}
		}
	}
}

// GetTurnState gets the current turn state
func (ga *GameAdapter) GetTurnState() interface{} {
	if ga.gameRef != nil {
		if gamePtr, ok := ga.gameRef.(*interface{}); ok {
			if game, ok := (*gamePtr).(interface{ GetTurnState() interface{} }); ok {
				return game.GetTurnState()
			}
		}
	}
	return nil
}

// GetTurnCounter gets the current turn counter
func (ga *GameAdapter) GetTurnCounter() int {
	if ga.gameRef != nil {
		if gamePtr, ok := ga.gameRef.(*interface{}); ok {
			if game, ok := (*gamePtr).(interface{ GetTurnCounter() int }); ok {
				return game.GetTurnCounter()
			}
		}
	}
	return 0
}

// SetTurnCounter sets the turn counter
func (ga *GameAdapter) SetTurnCounter(count int) {
	if ga.gameRef != nil {
		if gamePtr, ok := ga.gameRef.(*interface{}); ok {
			if game, ok := (*gamePtr).(interface{ SetTurnCounter(int) }); ok {
				game.SetTurnCounter(count)
			}
		}
	}
}

// IncrementTurnCounter increments the turn counter
func (ga *GameAdapter) IncrementTurnCounter() {
	if ga.gameRef != nil {
		if gamePtr, ok := ga.gameRef.(*interface{}); ok {
			if game, ok := (*gamePtr).(interface{ IncrementTurnCounter() }); ok {
				game.IncrementTurnCounter()
			}
		}
	}
}

// MapAdapter adapts the map interface for the MapSystem
type MapAdapter struct {
	mapRef *interface{} // Will be set later to avoid circular dependency
}

func (m *MapAdapter) UnblockTile(x, y int) {
	// Implementation will be added when we wire this up to the actual Game struct
}

func (m *MapAdapter) BlockTile(x, y int) {
	// Implementation will be added when we wire this up to the actual Game struct
}

func (m *MapAdapter) IsBlocked(x, y int) bool {
	return false // Placeholder
}
