package systems

import (
	"sync"
	"time"

	"github.com/caustin/rrogue/events"
	"github.com/caustin/rrogue/world"
)

// UIMessage represents a message in the UI system with metadata
type UIMessage struct {
	Text        string
	MessageType string
	Timestamp   time.Time
}

// UISystem handles user interface messages and display logic
type UISystem struct {
	world    world.WorldService
	eventBus *events.EventBus

	// Message management
	messages    []UIMessage
	mutex       sync.RWMutex
	maxMessages int
}

// NewUISystem creates a new UI system with dependencies
func NewUISystem(world world.WorldService, eventBus *events.EventBus) *UISystem {
	return &UISystem{
		world:       world,
		eventBus:    eventBus,
		messages:    make([]UIMessage, 0),
		maxMessages: 10, // Keep last 10 messages
	}
}

// RegisterHandlers subscribes the UI system to relevant events
func (ui *UISystem) RegisterHandlers() {
	ui.eventBus.Subscribe(events.MessageEventType, ui.HandleMessage)
	ui.eventBus.Subscribe(events.ClearMessagesEventType, ui.HandleClearMessages)
}

// HandleMessage processes message events and adds them to the message queue
func (ui *UISystem) HandleMessage(event events.Event) {
	messageEvent := event.(*events.MessageEvent)

	ui.mutex.Lock()
	defer ui.mutex.Unlock()

	// Add new message
	message := UIMessage{
		Text:        messageEvent.Message,
		MessageType: messageEvent.MessageType,
		Timestamp:   messageEvent.Timestamp(),
	}

	ui.messages = append(ui.messages, message)

	// Trim messages if we exceed max
	if len(ui.messages) > ui.maxMessages {
		ui.messages = ui.messages[1:]
	}
}

// HandleClearMessages processes clear message events
func (ui *UISystem) HandleClearMessages(event events.Event) {
	clearEvent := event.(*events.ClearMessagesEvent)

	ui.mutex.Lock()
	defer ui.mutex.Unlock()

	if clearEvent.ClearAll {
		ui.messages = make([]UIMessage, 0)
	} else {
		// Keep only recent messages (last 3)
		if len(ui.messages) > 3 {
			ui.messages = ui.messages[len(ui.messages)-3:]
		}
	}
}

// GetCurrentMessages returns the current messages for display
func (ui *UISystem) GetCurrentMessages() []UIMessage {
	ui.mutex.RLock()
	defer ui.mutex.RUnlock()

	// Return a copy to avoid race conditions
	result := make([]UIMessage, len(ui.messages))
	copy(result, ui.messages)
	return result
}

// GetMessageTexts returns just the text of current messages (for compatibility)
// Messages are returned in reverse order (latest first) for better UX
func (ui *UISystem) GetMessageTexts() []string {
	ui.mutex.RLock()
	defer ui.mutex.RUnlock()

	texts := make([]string, len(ui.messages))
	// Reverse order - latest messages first
	for i, msg := range ui.messages {
		texts[len(ui.messages)-1-i] = msg.Text
	}
	return texts
}

// AddMessage is a direct method to add messages (for compatibility during transition)
func (ui *UISystem) AddMessage(text, messageType string) {
	messageEvent := events.NewMessageEvent(text, messageType)
	ui.eventBus.Publish(messageEvent)
}

// ClearMessages triggers a clear messages event
func (ui *UISystem) ClearMessages(clearAll bool) {
	clearEvent := events.NewClearMessagesEvent(clearAll)
	ui.eventBus.Publish(clearEvent)
}

// GetMessageCount returns the current number of messages
func (ui *UISystem) GetMessageCount() int {
	ui.mutex.RLock()
	defer ui.mutex.RUnlock()
	return len(ui.messages)
}
