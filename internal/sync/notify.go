package sync

import (
	"fmt"
	"io"
	"os"
	"time"
)

// NotifyLevel controls which events trigger notifications.
type NotifyLevel int

const (
	NotifyNone    NotifyLevel = iota
	NotifyFailure             // only failures
	NotifyAll                 // failures and successes
)

// NotifyEvent holds data about a completed profile sync.
type NotifyEvent struct {
	Profile   string
	Success   bool
	Changes   int
	Duration  time.Duration
	Err       error
	Timestamp time.Time
}

// Notifier writes human-readable sync notifications to an output sink.
type Notifier struct {
	out   io.Writer
	level NotifyLevel
}

// NewNotifier creates a Notifier. If out is nil, os.Stderr is used.
func NewNotifier(out io.Writer, level NotifyLevel) *Notifier {
	if out == nil {
		out = os.Stderr
	}
	return &Notifier{out: out, level: level}
}

// Notify emits a notification for the given event if the level permits it.
func (n *Notifier) Notify(ev NotifyEvent) {
	if n.level == NotifyNone {
		return
	}
	if n.level == NotifyFailure && ev.Success {
		return
	}

	ts := ev.Timestamp
	if ts.IsZero() {
		ts = time.Now()
	}

	if ev.Success {
		fmt.Fprintf(n.out, "[%s] ✔ profile %q synced: %d change(s) in %s\n",
			ts.Format(time.RFC3339), ev.Profile, ev.Changes, ev.Duration.Round(time.Millisecond))
	} else {
		fmt.Fprintf(n.out, "[%s] ✘ profile %q failed after %s: %v\n",
			ts.Format(time.RFC3339), ev.Profile, ev.Duration.Round(time.Millisecond), ev.Err)
	}
}
