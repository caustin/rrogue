package events

import (
	"sync"
)

// EventHandler is a function that processes events
type EventHandler func(event Event)

// EventBus manages event subscription and publishing
type EventBus struct {
	subscribers map[EventType][]EventHandler
	mutex       sync.RWMutex
}

// NewEventBus creates a new event bus
func NewEventBus() *EventBus {
	return &EventBus{
		subscribers: make(map[EventType][]EventHandler),
	}
}

// Subscribe adds a handler for a specific event type
func (bus *EventBus) Subscribe(eventType EventType, handler EventHandler) {
	bus.mutex.Lock()
	defer bus.mutex.Unlock()

	bus.subscribers[eventType] = append(bus.subscribers[eventType], handler)
}

// Publish sends an event to all subscribed handlers
func (bus *EventBus) Publish(event Event) {
	bus.mutex.RLock()
	handlers := bus.subscribers[event.Type()]
	bus.mutex.RUnlock()

	// Call all handlers for this event type
	for _, handler := range handlers {
		handler(event)
	}
}

// PublishMany sends multiple events in sequence
func (bus *EventBus) PublishMany(events []Event) {
	for _, event := range events {
		bus.Publish(event)
	}
}

// Unsubscribe removes a handler (for cleanup)
// Note: This is a simple implementation that removes ALL handlers of a type
func (bus *EventBus) Unsubscribe(eventType EventType) {
	bus.mutex.Lock()
	defer bus.mutex.Unlock()

	delete(bus.subscribers, eventType)
}

// GetSubscriberCount returns the number of handlers for an event type (useful for testing)
func (bus *EventBus) GetSubscriberCount(eventType EventType) int {
	bus.mutex.RLock()
	defer bus.mutex.RUnlock()

	return len(bus.subscribers[eventType])
}
