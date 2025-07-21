package systems

import (
	"sync"

	"github.com/caustin/rrogue/events"
	"github.com/caustin/rrogue/world"
)

// TurnState represents the current state of the game turn
type TurnState int

const (
	WaitingForPlayerInput TurnState = iota
	ProcessingPlayerAction
	ProcessingMonsterTurn
	GameOver
)

// GameStateSystem handles game state transitions and game over conditions
type GameStateSystem struct {
	world    world.WorldService
	eventBus *events.EventBus

	// Direct references to Game struct fields (for migration phase)
	turnStateRef   interface{} // Generic interface to avoid import cycle
	turnCounterRef *int

	// Internal state tracking
	currentState TurnState
	turnCounter  int
	mutex        sync.RWMutex
}

// NewGameStateSystem creates a new game state system
func NewGameStateSystem(world world.WorldService, eventBus *events.EventBus) *GameStateSystem {
	return &GameStateSystem{
		world:        world,
		eventBus:     eventBus,
		currentState: WaitingForPlayerInput,
		turnCounter:  0,
	}
}

// SetGameReferences sets direct references to Game struct fields
func (gs *GameStateSystem) SetGameReferences(turnState interface{}, turnCounter *int) {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()

	gs.turnStateRef = turnState
	gs.turnCounterRef = turnCounter

	// Sync initial state from Game struct (using reflection-like approach)
	if turnCounter != nil {
		gs.turnCounter = *turnCounter
	}
}

// RegisterHandlers subscribes the game state system to relevant events
func (gs *GameStateSystem) RegisterHandlers() {
	gs.eventBus.Subscribe(events.DeathEventType, gs.HandleDeath)
	gs.eventBus.Subscribe(events.GameOverEventType, gs.HandleGameOver)
	gs.eventBus.Subscribe(events.TurnChangeEventType, gs.HandleTurnChange)
	gs.eventBus.Subscribe(events.TurnCounterEventType, gs.HandleTurnCounter)
}

// HandleDeath processes death events and manages game over conditions
func (gs *GameStateSystem) HandleDeath(event events.Event) {
	deathEvent := event.(*events.DeathEvent)

	if deathEvent.IsPlayer {
		// Player died - trigger game over
		gs.TriggerGameOver("player_death")
	}
	// Monster deaths don't affect game state directly
}

// HandleGameOver processes game over events
func (gs *GameStateSystem) HandleGameOver(event events.Event) {
	_ = event.(*events.GameOverEvent) // Extract but don't need to use

	gs.mutex.Lock()
	defer gs.mutex.Unlock()

	// Set game over state
	gs.currentState = GameOver

	// Sync to Game struct if references are available
	gs.syncTurnStateToGame(GameOver)

	// Publish turn change event for other systems
	turnChangeEvent := events.NewTurnChangeEvent(
		gs.turnStateToString(gs.currentState),
		"GameOver",
		gs.turnCounter,
	)
	gs.eventBus.Publish(turnChangeEvent)
}

// HandleTurnChange processes turn change events
func (gs *GameStateSystem) HandleTurnChange(event events.Event) {
	turnChangeEvent := event.(*events.TurnChangeEvent)

	gs.mutex.Lock()
	defer gs.mutex.Unlock()

	// Update internal state
	newState := gs.stringToTurnState(turnChangeEvent.ToState)
	gs.currentState = newState

	// Sync to Game struct if references are available
	gs.syncTurnStateToGame(newState)
}

// HandleTurnCounter processes turn counter events
func (gs *GameStateSystem) HandleTurnCounter(event events.Event) {
	counterEvent := event.(*events.TurnCounterEvent)

	gs.mutex.Lock()
	defer gs.mutex.Unlock()

	gs.turnCounter = counterEvent.TurnCount

	// Sync to Game struct if references are available
	if gs.turnCounterRef != nil {
		*gs.turnCounterRef = counterEvent.TurnCount
	}
}

// TriggerGameOver publishes a game over event
func (gs *GameStateSystem) TriggerGameOver(reason string) {
	gs.mutex.RLock()
	finalTurn := gs.turnCounter
	gs.mutex.RUnlock()

	gameOverEvent := events.NewGameOverEvent(reason, finalTurn)
	gs.eventBus.Publish(gameOverEvent)
}

// ChangeTurn publishes a turn change event
func (gs *GameStateSystem) ChangeTurn(toState TurnState) {
	gs.mutex.RLock()
	fromState := gs.currentState
	turnCount := gs.turnCounter
	gs.mutex.RUnlock()

	turnChangeEvent := events.NewTurnChangeEvent(
		gs.turnStateToString(fromState),
		gs.turnStateToString(toState),
		turnCount,
	)
	gs.eventBus.Publish(turnChangeEvent)
}

// IncrementTurn increments the turn counter and publishes event
func (gs *GameStateSystem) IncrementTurn() {
	gs.mutex.Lock()
	gs.turnCounter++
	newCount := gs.turnCounter
	gs.mutex.Unlock()

	counterEvent := events.NewTurnCounterEvent(newCount, 1)
	gs.eventBus.Publish(counterEvent)
}

// GetCurrentState returns the current turn state (thread-safe)
func (gs *GameStateSystem) GetCurrentState() TurnState {
	gs.mutex.RLock()
	defer gs.mutex.RUnlock()
	return gs.currentState
}

// GetTurnCounter returns the current turn counter (thread-safe)
func (gs *GameStateSystem) GetTurnCounter() int {
	gs.mutex.RLock()
	defer gs.mutex.RUnlock()
	return gs.turnCounter
}

// syncTurnStateToGame syncs the internal turn state to the Game struct
func (gs *GameStateSystem) syncTurnStateToGame(state TurnState) {
	// Note: During migration, we need to convert our TurnState to game.TurnState
	// Since both use the same underlying int values, we can do a simple cast
	if gs.turnStateRef != nil {
		// Direct pointer assignment since TurnState values match game.TurnState values
		if ptr, ok := gs.turnStateRef.(*int); ok {
			*ptr = int(state)
		}
	}
}

// Helper functions for state conversion
func (gs *GameStateSystem) turnStateToString(state TurnState) string {
	switch state {
	case WaitingForPlayerInput:
		return "WaitingForPlayerInput"
	case ProcessingPlayerAction:
		return "ProcessingPlayerAction"
	case ProcessingMonsterTurn:
		return "ProcessingMonsterTurn"
	case GameOver:
		return "GameOver"
	default:
		return "Unknown"
	}
}

func (gs *GameStateSystem) stringToTurnState(state string) TurnState {
	switch state {
	case "WaitingForPlayerInput":
		return WaitingForPlayerInput
	case "ProcessingPlayerAction":
		return ProcessingPlayerAction
	case "ProcessingMonsterTurn":
		return ProcessingMonsterTurn
	case "GameOver":
		return GameOver
	default:
		return WaitingForPlayerInput
	}
}
