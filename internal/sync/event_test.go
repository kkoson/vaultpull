package sync

import (
	"errors"
	"testing"
	"time"
)

func TestNewEventBus_Empty(t *testing.T) {
	bus := NewEventBus()
	if bus == nil {
		t.Fatal("expected non-nil EventBus")
	}
	if bus.SubscriberCount() != 0 {
		t.Fatalf("expected 0 subscribers, got %d", bus.SubscriberCount())
	}
}

func TestEventBus_Subscribe_NilIgnored(t *testing.T) {
	bus := NewEventBus()
	bus.Subscribe(nil)
	if bus.SubscriberCount() != 0 {
		t.Fatalf("expected 0 subscribers after nil, got %d", bus.SubscriberCount())
	}
}

func TestEventBus_Subscribe_IncrementsCount(t *testing.T) {
	bus := NewEventBus()
	bus.Subscribe(func(Event) {})
	bus.Subscribe(func(Event) {})
	if bus.SubscriberCount() != 2 {
		t.Fatalf("expected 2 subscribers, got %d", bus.SubscriberCount())
	}
}

func TestEventBus_Publish_DeliveredToAll(t *testing.T) {
	bus := NewEventBus()
	var received []Event
	bus.Subscribe(func(e Event) { received = append(received, e) })
	bus.Subscribe(func(e Event) { received = append(received, e) })

	bus.Publish(Event{Type: EventProfileStarted, Profile: "prod"})

	if len(received) != 2 {
		t.Fatalf("expected 2 deliveries, got %d", len(received))
	}
	for _, e := range received {
		if e.Profile != "prod" {
			t.Errorf("unexpected profile %q", e.Profile)
		}
	}
}

func TestEventBus_Publish_SetsTimestamp(t *testing.T) {
	bus := NewEventBus()
	var got Event
	bus.Subscribe(func(e Event) { got = e })

	before := time.Now()
	bus.Publish(Event{Type: EventSyncStarted})
	after := time.Now()

	if got.Timestamp.Before(before) || got.Timestamp.After(after) {
		t.Errorf("timestamp %v not in expected range [%v, %v]", got.Timestamp, before, after)
	}
}

func TestEventBus_Publish_PreservesExistingTimestamp(t *testing.T) {
	bus := NewEventBus()
	var got Event
	bus.Subscribe(func(e Event) { got = e })

	fixed := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	bus.Publish(Event{Type: EventSyncCompleted, Timestamp: fixed})

	if !got.Timestamp.Equal(fixed) {
		t.Errorf("expected timestamp %v, got %v", fixed, got.Timestamp)
	}
}

func TestEventBus_Publish_WithError(t *testing.T) {
	bus := NewEventBus()
	var got Event
	bus.Subscribe(func(e Event) { got = e })

	sentinel := errors.New("vault unavailable")
	bus.Publish(Event{Type: EventProfileFailed, Profile: "staging", Err: sentinel})

	if got.Err == nil || got.Err.Error() != sentinel.Error() {
		t.Errorf("expected error %v, got %v", sentinel, got.Err)
	}
}
