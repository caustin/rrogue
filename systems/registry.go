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

	// For now, create placeholder systems - we'll wire them up properly later
	// registry.GameState = NewGameStateSystem(eventBus, world, &GameStateAdapter{})
	// registry.Map = NewMapSystem(eventBus, world, &MapAdapter{})

	return registry
}

// RegisterAllHandlers subscribes all systems to their respective events
func (r *SystemRegistry) RegisterAllHandlers() {
	r.Combat.RegisterHandlers()

	// Register other systems when they're implemented
	// if r.GameState != nil {
	//     r.GameState.RegisterHandlers()
	// }
	// if r.Map != nil {
	//     r.Map.RegisterHandlers()
	// }
}

// GameStateAdapter adapts the game state interface for the GameStateSystem
type GameStateAdapter struct {
	gameRef *interface{} // Will be set later to avoid circular dependency
}

func (g *GameStateAdapter) SetGameOver() {
	// Implementation will be added when we wire this up to the actual Game struct
}

func (g *GameStateAdapter) GetCurrentTurn() string {
	return "unknown" // Placeholder
}

func (g *GameStateAdapter) SetCurrentTurn(turn string) {
	// Implementation will be added when we wire this up to the actual Game struct
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
