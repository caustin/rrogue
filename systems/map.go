package systems

import (
	"github.com/caustin/rrogue/events"
	"github.com/caustin/rrogue/world"
)

// MapSystem handles map-related operations
type MapSystem struct {
	eventBus   *events.EventBus
	world      world.WorldService
	mapManager MapManager
}

// MapManager interface for map operations
type MapManager interface {
	UnblockTile(x, y int)
	BlockTile(x, y int)
	IsBlocked(x, y int) bool
}

// NewMapSystem creates a new map system
func NewMapSystem(eventBus *events.EventBus, world world.WorldService, mapManager MapManager) *MapSystem {
	return &MapSystem{
		eventBus:   eventBus,
		world:      world,
		mapManager: mapManager,
	}
}

// RegisterHandlers subscribes the map system to relevant events
func (ms *MapSystem) RegisterHandlers() {
	ms.eventBus.Subscribe(events.DeathEventType, ms.HandleEntityDeath)
	ms.eventBus.Subscribe(events.MoveEventType, ms.HandleEntityMove)
	ms.eventBus.Subscribe(events.TileBlockedEventType, ms.HandleTileBlocked)
	ms.eventBus.Subscribe(events.TileUnblockedEventType, ms.HandleTileUnblocked)
}

// HandleEntityDeath processes entity death and unblocks tiles
func (ms *MapSystem) HandleEntityDeath(event events.Event) {
	deathEvent := event.(*events.DeathEvent)

	if !deathEvent.IsPlayer {
		// Monster died - unblock the tile
		ms.mapManager.UnblockTile(deathEvent.Position.X, deathEvent.Position.Y)

		// Dispose the entity from the world
		ms.world.DisposeEntity(deathEvent.Entity)

		// Publish tile unblocked event
		tileEvent := events.NewTileUnblockedEvent(deathEvent.Position, "monster_death")
		ms.eventBus.Publish(tileEvent)
	}
}

// HandleEntityMove processes entity movement and updates tile blocking
func (ms *MapSystem) HandleEntityMove(event events.Event) {
	moveEvent := event.(*events.MoveEvent)

	// Unblock old position
	ms.mapManager.UnblockTile(moveEvent.FromPos.X, moveEvent.FromPos.Y)

	// Block new position
	ms.mapManager.BlockTile(moveEvent.ToPos.X, moveEvent.ToPos.Y)

	// Publish tile events
	unblockEvent := events.NewTileUnblockedEvent(moveEvent.FromPos, "entity_move")
	blockEvent := events.NewTileBlockedEvent(moveEvent.ToPos, "entity_move")

	ms.eventBus.Publish(unblockEvent)
	ms.eventBus.Publish(blockEvent)
}

// HandleTileBlocked processes tile blocked events
func (ms *MapSystem) HandleTileBlocked(event events.Event) {
	tileEvent := event.(*events.TileBlockedEvent)
	ms.mapManager.BlockTile(tileEvent.Position.X, tileEvent.Position.Y)
}

// HandleTileUnblocked processes tile unblocked events
func (ms *MapSystem) HandleTileUnblocked(event events.Event) {
	tileEvent := event.(*events.TileUnblockedEvent)
	ms.mapManager.UnblockTile(tileEvent.Position.X, tileEvent.Position.Y)
}
