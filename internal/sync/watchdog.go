package sync

import (
	"context"
	"sync"
	"time"
)

// Watchdog monitors a set of named goroutines and restarts them if they
// stop reporting within a configurable deadline.
type Watchdog struct {
	mu       sync.Mutex
	heartbeat map[string]time.Time
	deadline  time.Duration
	restart   func(name string)
}

// NewWatchdog creates a Watchdog that calls restart(name) whenever a
// registered worker has not sent a heartbeat within deadline.
// A zero or negative deadline defaults to 30 seconds.
func NewWatchdog(deadline time.Duration, restart func(name string)) *Watchdog {
	if deadline <= 0 {
		deadline = 30 * time.Second
	}
	if restart == nil {
		restart = func(string) {}
	}
	return &Watchdog{
		heartbeat: make(map[string]time.Time),
		deadline:  deadline,
		restart:   restart,
	}
}

// Beat records a heartbeat for the named worker. If the worker is not yet
// registered, this call also registers it.
func (w *Watchdog) Beat(name string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.heartbeat[name] = time.Now()
}

// Unregister removes a worker from monitoring.
func (w *Watchdog) Unregister(name string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	delete(w.heartbeat, name)
}

// Registered returns the names of all currently monitored workers.
func (w *Watchdog) Registered() []string {
	w.mu.Lock()
	defer w.mu.Unlock()
	names := make([]string, 0, len(w.heartbeat))
	for name := range w.heartbeat {
		names = append(names, name)
	}
	return names
}

// Start begins the watchdog loop, checking heartbeats every interval until
// ctx is cancelled. interval defaults to deadline/2 when zero.
func (w *Watchdog) Start(ctx context.Context, interval time.Duration) {
	if interval <= 0 {
		interval = w.deadline / 2
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case now := <-ticker.C:
			w.check(now)
		}
	}
}

func (w *Watchdog) check(now time.Time) {
	w.mu.Lock()
	var stale []string
	for name, last := range w.heartbeat {
		if now.Sub(last) > w.deadline {
			stale = append(stale, name)
		}
	}
	w.mu.Unlock()
	for _, name := range stale {
		w.restart(name)
	}
}
