package sync

import (
	"context"
	"errors"
)

// Limiter caps the number of concurrent profile sync operations.
// It is backed by a buffered channel used as a counting semaphore.
type Limiter struct {
	sem chan struct{}
}

// NewLimiter returns a Limiter that allows at most n concurrent
// operations. If n is less than 1 it is clamped to 1.
func NewLimiter(n int) *Limiter {
	if n < 1 {
		n = 1
	}
	return &Limiter{sem: make(chan struct{}, n)}
}

// Run acquires a slot, executes fn, then releases the slot.
// If ctx is cancelled before a slot is available, Run returns
// ctx.Err() without calling fn.
func (l *Limiter) Run(ctx context.Context, fn func() error) error {
	select {
	case l.sem <- struct{}{}:
		defer func() { <-l.sem }()
		return fn()
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Available returns the number of additional slots that can be
// acquired without blocking.
func (l *Limiter) Available() int {
	return cap(l.sem) - len(l.sem)
}

// Capacity returns the maximum number of concurrent operations
// allowed by this Limiter.
func (l *Limiter) Capacity() int {
	return cap(l.sem)
}

// ErrLimiterNilFn is returned when Run is called with a nil function.
var ErrLimiterNilFn = errors.New("limiter: fn must not be nil")
