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
	Reason string // "player_death", "victory", etc.
}

func NewGameOverEvent(reason string) *GameOverEvent {
	return &GameOverEvent{
		BaseEvent: NewBaseEvent(GameOverEventType),
		Reason:    reason,
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
