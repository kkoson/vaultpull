package sync

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// TraceLevel controls verbosity of trace output.
type TraceLevel int

const (
	TraceLevelOff   TraceLevel = iota
	TraceLevelBasic            // profile start/end
	TraceLevelFull             // includes stage timings
)

// TraceEntry records a single trace event.
type TraceEntry struct {
	Profile   string
	Stage     string
	Event     string
	Duration  time.Duration
	Timestamp time.Time
}

// Tracer collects and emits trace events for sync operations.
type Tracer struct {
	mu      sync.Mutex
	level   TraceLevel
	w       io.Writer
	entries []TraceEntry
}

// NewTracer creates a Tracer at the given level. If w is nil, os.Stderr is used.
func NewTracer(level TraceLevel, w io.Writer) *Tracer {
	if w == nil {
		w = os.Stderr
	}
	return &Tracer{level: level, w: w}
}

// Record adds a trace entry if the tracer level is sufficient.
func (t *Tracer) Record(profile, stage, event string, d time.Duration) {
	if t == nil || t.level == TraceLevelOff {
		return
	}
	if t.level == TraceLevelBasic && stage != "" {
		return
	}
	e := TraceEntry{
		Profile:   profile,
		Stage:     stage,
		Event:     event,
		Duration:  d,
		Timestamp: time.Now().UTC(),
	}
	t.mu.Lock()
	t.entries = append(t.entries, e)
	t.mu.Unlock()
	t.emit(e)
}

// Entries returns a snapshot of all recorded entries.
func (t *Tracer) Entries() []TraceEntry {
	t.mu.Lock()
	defer t.mu.Unlock()
	out := make([]TraceEntry, len(t.entries))
	copy(out, t.entries)
	return out
}

func (t *Tracer) emit(e TraceEntry) {
	loc := e.Profile
	if e.Stage != "" {
		loc = fmt.Sprintf("%s/%s", e.Profile, e.Stage)
	}
	var dur string
	if e.Duration > 0 {
		dur = fmt.Sprintf(" (%s)", e.Duration.Round(time.Millisecond))
	}
	fmt.Fprintf(t.w, "[trace] %s %s%s\n", loc, e.Event, dur)
}
