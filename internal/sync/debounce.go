package sync

import (
	"context"
	"sync"
	"time"
)

// Debouncer delays execution of a function until after a quiet period has elapsed.
// If the function is triggered again before the timer fires, the timer resets.
// This is useful for batching rapid profile sync requests into a single run.
type Debouncer struct {
	delay  time.Duration
	mu     sync.Mutex
	timer  *time.Timer
	cancel context.CancelFunc
}

// NewDebouncer creates a Debouncer with the given delay.
// If delay is zero or negative, a default of 500ms is used.
func NewDebouncer(delay time.Duration) *Debouncer {
	if delay <= 0 {
		delay = 500 * time.Millisecond
	}
	return &Debouncer{delay: delay}
}

// Trigger schedules fn to run after the debounce delay.
// If Trigger is called again before the delay elapses, the previous
// scheduled call is cancelled and the timer resets.
// The provided context governs the lifetime of the debounced call.
func (d *Debouncer) Trigger(ctx context.Context, fn func()) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.timer != nil {
		d.timer.Stop()
	}
	if d.cancel != nil {
		d.cancel()
	}

	ctx, cancel := context.WithCancel(ctx)
	d.cancel = cancel

	d.timer = time.AfterFunc(d.delay, func() {
		select {
		case <-ctx.Done():
			return
		default:
			fn()
		}
	})
}

// Flush cancels any pending debounced call without executing it.
func (d *Debouncer) Flush() {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.timer != nil {
		d.timer.Stop()
		d.timer = nil
	}
	if d.cancel != nil {
		d.cancel()
		d.cancel = nil
	}
}
