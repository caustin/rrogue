package systems

import (
	"github.com/caustin/rrogue/events"
	"github.com/caustin/rrogue/world"
)

// GameStateSystem handles game state transitions and game over conditions
type GameStateSystem struct {
	eventBus  *events.EventBus
	world     world.WorldService
	gameState GameStateManager
}

// GameStateManager interface for managing game state
type GameStateManager interface {
	SetGameOver()
	GetCurrentTurn() string
	SetCurrentTurn(turn string)
}

// NewGameStateSystem creates a new game state system
func NewGameStateSystem(eventBus *events.EventBus, world world.WorldService, gameState GameStateManager) *GameStateSystem {
	return &GameStateSystem{
		eventBus:  eventBus,
		world:     world,
		gameState: gameState,
	}
}

// RegisterHandlers subscribes the game state system to relevant events
func (gs *GameStateSystem) RegisterHandlers() {
	gs.eventBus.Subscribe(events.DeathEventType, gs.HandleDeath)
	gs.eventBus.Subscribe(events.TurnStartEventType, gs.HandleTurnStart)
	gs.eventBus.Subscribe(events.TurnEndEventType, gs.HandleTurnEnd)
}

// HandleDeath processes entity death events and checks for game over
func (gs *GameStateSystem) HandleDeath(event events.Event) {
	deathEvent := event.(*events.DeathEvent)

	if deathEvent.IsPlayer {
		// Player died - game over
		defenderMessage := gs.world.GetUserMessage(deathEvent.Entity)
		defenderMessage.GameStateMessage = "Game Over!\n"

		// Publish game over event
		gameOverEvent := events.NewGameOverEvent("player_death")
		gs.eventBus.Publish(gameOverEvent)

		// Set game state to game over
		gs.gameState.SetGameOver()
	}
}

// HandleTurnStart processes turn start events
func (gs *GameStateSystem) HandleTurnStart(event events.Event) {
	turnEvent := event.(*events.TurnStartEvent)
	gs.gameState.SetCurrentTurn(turnEvent.TurnType)
}

// HandleTurnEnd processes turn end events
func (gs *GameStateSystem) HandleTurnEnd(event events.Event) {
	// Additional turn end logic can go here
}
