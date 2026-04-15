package sync

import (
	"io"
	"os"
	"sync"
	"time"
)

// ObserveLevel controls how much detail the observer emits.
type ObserveLevel int

const (
	// ObserveOff disables all observation output.
	ObserveOff ObserveLevel = iota
	// ObserveSummary emits only per-profile outcome lines.
	ObserveSummary
	// ObserveFull emits per-key change details as well.
	ObserveFull
)

// ObserveEvent holds a single observation record.
type ObserveEvent struct {
	Profile   string
	Key       string
	Change    string // "added", "updated", "removed", "unchanged"
	Timestamp time.Time
}

// Observer collects and optionally prints sync observations.
type Observer struct {
	mu     sync.Mutex
	level  ObserveLevel
	out    io.Writer
	events []ObserveEvent
}

// NewObserver creates an Observer at the given level.
// If w is nil, os.Stderr is used.
func NewObserver(level ObserveLevel, w io.Writer) *Observer {
	if w == nil {
		w = os.Stderr
	}
	return &Observer{level: level, out: w}
}

// Record stores an observation event and optionally prints it.
func (o *Observer) Record(profile, key, change string) {
	if o == nil || o.level == ObserveOff {
		return
	}
	ev := ObserveEvent{
		Profile:   profile,
		Key:       key,
		Change:    change,
		Timestamp: time.Now().UTC(),
	}
	o.mu.Lock()
	o.events = append(o.events, ev)
	o.mu.Unlock()

	if o.level >= ObserveFull && key != "" {
		io.WriteString(o.out, "[observe] "+profile+" "+key+" "+change+"\n") //nolint:errcheck
	} else if o.level == ObserveSummary && key == "" {
		io.WriteString(o.out, "[observe] "+profile+" "+change+"\n") //nolint:errcheck
	}
}

// Events returns a snapshot of all collected events.
func (o *Observer) Events() []ObserveEvent {
	if o == nil {
		return nil
	}
	o.mu.Lock()
	defer o.mu.Unlock()
	out := make([]ObserveEvent, len(o.events))
	copy(out, o.events)
	return out
}

// Reset clears all recorded events.
func (o *Observer) Reset() {
	if o == nil {
		return
	}
	o.mu.Lock()
	o.events = o.events[:0]
	o.mu.Unlock()
}
