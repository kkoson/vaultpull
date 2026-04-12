package sync

import (
	"context"
	"sync"
	"time"
)

// ThrottleConfig controls how profile syncs are throttled over time.
type ThrottleConfig struct {
	// MinInterval is the minimum time between successive sync operations.
	MinInterval time.Duration
	// BurstSize is the number of syncs allowed before throttling kicks in.
	BurstSize int
}

// DefaultThrottleConfig returns a ThrottleConfig with sensible defaults.
func DefaultThrottleConfig() ThrottleConfig {
	return ThrottleConfig{
		MinInterval: 200 * time.Millisecond,
		BurstSize:   3,
	}
}

// Throttle enforces a minimum interval between sync operations with burst support.
type Throttle struct {
	cfg     ThrottleConfig
	mu      sync.Mutex
	last    time.Time
	bursts  int
}

// NewThrottle creates a Throttle from the given config.
func NewThrottle(cfg ThrottleConfig) *Throttle {
	if cfg.MinInterval <= 0 {
		cfg.MinInterval = DefaultThrottleConfig().MinInterval
	}
	if cfg.BurstSize <= 0 {
		cfg.BurstSize = 1
	}
	return &Throttle{cfg: cfg}
}

// Wait blocks until the throttle allows the next operation, or until ctx is
// cancelled. Returns ctx.Err() if the context is done while waiting.
func (t *Throttle) Wait(ctx context.Context) error {
	t.mu.Lock()
	now := time.Now()
	var delay time.Duration
	if !t.last.IsZero() && t.bursts >= t.cfg.BurstSize {
		elapsed := now.Sub(t.last)
		if elapsed < t.cfg.MinInterval {
			delay = t.cfg.MinInterval - elapsed
		}
	}
	if delay == 0 {
		if t.last.IsZero() || now.Sub(t.last) >= t.cfg.MinInterval {
			t.bursts = 0
		}
		t.bursts++
		t.last = now
		t.mu.Unlock()
		return nil
	}
	t.mu.Unlock()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(delay):
		t.mu.Lock()
		t.bursts = 1
		t.last = time.Now()
		t.mu.Unlock()
		return nil
	}
}

// Reset clears the throttle state.
func (t *Throttle) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.last = time.Time{}
	t.bursts = 0
}
