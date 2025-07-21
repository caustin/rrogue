package events

// UI Event Types
const (
	MessageEventType       EventType = "message_event"
	ClearMessagesEventType EventType = "clear_messages_event"
	UIUpdateEventType      EventType = "ui_update_event"
)

// MessageEvent represents a message to be displayed in the UI log
type MessageEvent struct {
	BaseEvent
	Message     string
	MessageType string // "attack", "death", "gamestate", "info", etc.
}

// NewMessageEvent creates a new message event
func NewMessageEvent(message, messageType string) *MessageEvent {
	return &MessageEvent{
		BaseEvent:   NewBaseEvent(MessageEventType),
		Message:     message,
		MessageType: messageType,
	}
}

// Type returns the event type
func (e *MessageEvent) Type() EventType {
	return MessageEventType
}

// ClearMessagesEvent represents a request to clear the message log
type ClearMessagesEvent struct {
	BaseEvent
	ClearAll bool // If true, clear all messages; if false, clear old messages
}

// NewClearMessagesEvent creates a new clear messages event
func NewClearMessagesEvent(clearAll bool) *ClearMessagesEvent {
	return &ClearMessagesEvent{
		BaseEvent: NewBaseEvent(ClearMessagesEventType),
		ClearAll:  clearAll,
	}
}

// Type returns the event type
func (e *ClearMessagesEvent) Type() EventType {
	return ClearMessagesEventType
}

// UIUpdateEvent represents a request to update the UI
type UIUpdateEvent struct {
	BaseEvent
	UpdateType string // "log", "hud", "all", etc.
}

// NewUIUpdateEvent creates a new UI update event
func NewUIUpdateEvent(updateType string) *UIUpdateEvent {
	return &UIUpdateEvent{
		BaseEvent:  NewBaseEvent(UIUpdateEventType),
		UpdateType: updateType,
	}
}

// Type returns the event type
func (e *UIUpdateEvent) Type() EventType {
	return UIUpdateEventType
}
