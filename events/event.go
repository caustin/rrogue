package events

import (
	"time"
)

// Event represents something that happened in the game
type Event interface {
	Type() EventType
	Timestamp() time.Time
}

// EventType identifies different kinds of events
type EventType string

const (
	// Combat Events
	AttackEventType EventType = "attack"
	DamageEventType EventType = "damage"
	DeathEventType  EventType = "death"

	// Movement Events
	MoveEventType EventType = "move"

	// Game State Events
	TurnStartEventType   EventType = "turn_start"
	TurnEndEventType     EventType = "turn_end"
	TurnChangeEventType  EventType = "turn_change"
	TurnCounterEventType EventType = "turn_counter"
	GameOverEventType    EventType = "game_over"

	// Map Events
	TileBlockedEventType   EventType = "tile_blocked"
	TileUnblockedEventType EventType = "tile_unblocked"
)

// BaseEvent provides common event functionality
type BaseEvent struct {
	EventType EventType
	Time      time.Time
}

func NewBaseEvent(eventType EventType) BaseEvent {
	return BaseEvent{
		EventType: eventType,
		Time:      time.Now(),
	}
}

func (e BaseEvent) Type() EventType {
	return e.EventType
}

func (e BaseEvent) Timestamp() time.Time {
	return e.Time
}
