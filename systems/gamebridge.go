package systems

import (
	"github.com/caustin/rrogue/events"
)

// GameBridge provides a way for systems to interact with game state
// This is a temporary solution until we have full event-driven game state management
type GameBridge struct {
	eventBus *events.EventBus
	gameRef  interface{} // Will hold reference to Game struct
}

// NewGameBridge creates a bridge between systems and game state
func NewGameBridge(eventBus *events.EventBus) *GameBridge {
	bridge := &GameBridge{
		eventBus: eventBus,
	}

	// Subscribe to death events to handle game over
	eventBus.Subscribe(events.DeathEventType, bridge.HandleDeath)

	return bridge
}

// SetGameReference allows the bridge to reference the main game
func (gb *GameBridge) SetGameReference(game interface{}) {
	gb.gameRef = game
}

// HandleDeath processes death events and manages game state changes
func (gb *GameBridge) HandleDeath(event events.Event) {
	deathEvent := event.(*events.DeathEvent)

	if deathEvent.IsPlayer {
		// Handle player death - set game over
		// We'll need to find a way to set the turn state
		// For now, this just ensures the event is properly handled
	} else {
		// Handle monster death - for map cleanup this will be handled
		// when we implement the proper MapSystem
	}
}
