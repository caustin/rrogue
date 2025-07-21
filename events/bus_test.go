package events

import (
	"github.com/caustin/rrogue/components"
	"testing"
	"time"
)

func TestEventBus_Subscribe_and_Publish(t *testing.T) {
	bus := NewEventBus()

	// Track if handler was called
	handlerCalled := false
	var receivedEvent Event

	// Subscribe to attack events
	bus.Subscribe(AttackEventType, func(event Event) {
		handlerCalled = true
		receivedEvent = event
	})

	// Create and publish an attack event
	attackerPos := &components.Position{X: 1, Y: 1}
	defenderPos := &components.Position{X: 2, Y: 2}
	attackEvent := NewAttackEvent(nil, nil, attackerPos, defenderPos, 5, true)
	bus.Publish(attackEvent)

	// Verify handler was called
	if !handlerCalled {
		t.Error("Event handler was not called")
	}

	// Verify correct event was received
	if receivedEvent == nil {
		t.Error("No event was received")
	}

	if receivedEvent.Type() != AttackEventType {
		t.Errorf("Expected event type %s, got %s", AttackEventType, receivedEvent.Type())
	}
}

func TestEventBus_MultipleSubscribers(t *testing.T) {
	bus := NewEventBus()

	// Track how many handlers were called
	callCount := 0

	// Subscribe multiple handlers to the same event type
	bus.Subscribe(DamageEventType, func(event Event) {
		callCount++
	})

	bus.Subscribe(DamageEventType, func(event Event) {
		callCount++
	})

	// Publish event
	damageEvent := NewDamageEvent(nil, 10, "sword", false)
	bus.Publish(damageEvent)

	// Verify both handlers were called
	if callCount != 2 {
		t.Errorf("Expected 2 handler calls, got %d", callCount)
	}
}

func TestEventBus_GetSubscriberCount(t *testing.T) {
	bus := NewEventBus()

	// Initially no subscribers
	count := bus.GetSubscriberCount(AttackEventType)
	if count != 0 {
		t.Errorf("Expected 0 subscribers, got %d", count)
	}

	// Add subscribers
	bus.Subscribe(AttackEventType, func(event Event) {})
	bus.Subscribe(AttackEventType, func(event Event) {})

	count = bus.GetSubscriberCount(AttackEventType)
	if count != 2 {
		t.Errorf("Expected 2 subscribers, got %d", count)
	}
}

func TestEventCreation(t *testing.T) {
	// Test that events are created with correct timestamps
	before := time.Now()
	event := NewGameOverEvent("player_death")
	after := time.Now()

	if event.Timestamp().Before(before) || event.Timestamp().After(after) {
		t.Error("Event timestamp is not within expected range")
	}

	if event.Type() != GameOverEventType {
		t.Errorf("Expected event type %s, got %s", GameOverEventType, event.Type())
	}
}
