package systems

import (
	"github.com/caustin/rrogue/events"
)

// MapBridge handles map-related events temporarily until we have a full MapSystem
type MapBridge struct {
	eventBus *events.EventBus
	gameRef  interface{} // Will hold reference to Game struct for map access
}

// NewMapBridge creates a bridge for map operations
func NewMapBridge(eventBus *events.EventBus) *MapBridge {
	bridge := &MapBridge{
		eventBus: eventBus,
	}

	// Subscribe to death events to handle tile unblocking
	eventBus.Subscribe(events.DeathEventType, bridge.HandleEntityDeath)

	return bridge
}

// SetGameReference allows the bridge to reference the main game for map access
func (mb *MapBridge) SetGameReference(game interface{}) {
	mb.gameRef = game
}

// HandleEntityDeath processes death events and unblocks tiles
func (mb *MapBridge) HandleEntityDeath(event events.Event) {
	deathEvent := event.(*events.DeathEvent)

	if !deathEvent.IsPlayer {
		// For monster death, we need to unblock the tile
		// This is a temporary hack until we have proper MapSystem
		// The actual tile unblocking will need to be done by the caller
		// since we can't access the map from here without circular dependencies
	}
}
