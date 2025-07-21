package events

import (
	"github.com/caustin/rrogue/components"
)

// TurnStartEvent represents the start of a turn
type TurnStartEvent struct {
	BaseEvent
	TurnType    string // "player" or "monster"
	TurnCounter int
}

func NewTurnStartEvent(turnType string, turnCounter int) *TurnStartEvent {
	return &TurnStartEvent{
		BaseEvent:   NewBaseEvent(TurnStartEventType),
		TurnType:    turnType,
		TurnCounter: turnCounter,
	}
}

// TurnEndEvent represents the end of a turn
type TurnEndEvent struct {
	BaseEvent
	TurnType    string
	TurnCounter int
}

func NewTurnEndEvent(turnType string, turnCounter int) *TurnEndEvent {
	return &TurnEndEvent{
		BaseEvent:   NewBaseEvent(TurnEndEventType),
		TurnType:    turnType,
		TurnCounter: turnCounter,
	}
}

// GameOverEvent represents the game ending
type GameOverEvent struct {
	BaseEvent
	Reason    string // "player_death", "victory", etc.
	FinalTurn int    // Turn number when game ended
}

func NewGameOverEvent(reason string, finalTurn int) *GameOverEvent {
	return &GameOverEvent{
		BaseEvent: NewBaseEvent(GameOverEventType),
		Reason:    reason,
		FinalTurn: finalTurn,
	}
}

// TurnChangeEvent represents a change in turn state
type TurnChangeEvent struct {
	BaseEvent
	FromState string
	ToState   string
	TurnCount int
}

// NewTurnChangeEvent creates a new turn change event
func NewTurnChangeEvent(fromState, toState string, turnCount int) *TurnChangeEvent {
	return &TurnChangeEvent{
		BaseEvent: NewBaseEvent(TurnChangeEventType),
		FromState: fromState,
		ToState:   toState,
		TurnCount: turnCount,
	}
}

// TurnCounterEvent represents turn counter updates
type TurnCounterEvent struct {
	BaseEvent
	TurnCount int
	Increment int
}

// NewTurnCounterEvent creates a new turn counter event
func NewTurnCounterEvent(turnCount, increment int) *TurnCounterEvent {
	return &TurnCounterEvent{
		BaseEvent: NewBaseEvent(TurnCounterEventType),
		TurnCount: turnCount,
		Increment: increment,
	}
}

// TileBlockedEvent represents a tile becoming blocked
type TileBlockedEvent struct {
	BaseEvent
	Position *components.Position
	Reason   string // "monster_spawn", "player_move", etc.
}

func NewTileBlockedEvent(position *components.Position, reason string) *TileBlockedEvent {
	return &TileBlockedEvent{
		BaseEvent: NewBaseEvent(TileBlockedEventType),
		Position:  position,
		Reason:    reason,
	}
}

// TileUnblockedEvent represents a tile becoming unblocked
type TileUnblockedEvent struct {
	BaseEvent
	Position *components.Position
	Reason   string // "monster_death", "player_move", etc.
}

func NewTileUnblockedEvent(position *components.Position, reason string) *TileUnblockedEvent {
	return &TileUnblockedEvent{
		BaseEvent: NewBaseEvent(TileUnblockedEventType),
		Position:  position,
		Reason:    reason,
	}
}
