package sync

import "time"

// EventType identifies the kind of sync lifecycle event.
type EventType string

const (
	EventProfileStarted  EventType = "profile.started"
	EventProfileSuccess  EventType = "profile.success"
	EventProfileFailed   EventType = "profile.failed"
	EventSyncStarted     EventType = "sync.started"
	EventSyncCompleted   EventType = "sync.completed"
)

// Event carries information about a point-in-time occurrence during a sync run.
type Event struct {
	Type      EventType
	Profile   string
	Timestamp time.Time
	Err       error
	Meta      map[string]string
}

// EventHandler is a function that receives a sync Event.
type EventHandler func(e Event)

// EventBus dispatches sync lifecycle events to registered handlers.
type EventBus struct {
	handlers []EventHandler
}

// NewEventBus returns an initialised EventBus.
func NewEventBus() *EventBus {
	return &EventBus{}
}

// Subscribe registers a handler that will be called for every published event.
func (b *EventBus) Subscribe(h EventHandler) {
	if h != nil {
		b.handlers = append(b.handlers, h)
	}
}

// Publish sends the event to all registered handlers.
func (b *EventBus) Publish(e Event) {
	if e.Timestamp.IsZero() {
		e.Timestamp = time.Now()
	}
	for _, h := range b.handlers {
		h(e)
	}
}

// SubscriberCount returns the number of registered handlers.
func (b *EventBus) SubscriberCount() int {
	return len(b.handlers)
}
