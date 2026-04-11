package sync

import (
	"context"
	"fmt"
)

// Semaphore limits the number of concurrent profile sync operations.
type Semaphore struct {
	ch chan struct{}
}

// NewSemaphore creates a Semaphore that allows at most n concurrent acquisitions.
// Returns an error if n is less than 1.
func NewSemaphore(n int) (*Semaphore, error) {
	if n < 1 {
		return nil, fmt.Errorf("semaphore: concurrency limit must be at least 1, got %d", n)
	}
	return &Semaphore{ch: make(chan struct{}, n)}, nil
}

// Acquire blocks until a slot is available or ctx is cancelled.
// Returns ctx.Err() if the context is done before a slot is acquired.
func (s *Semaphore) Acquire(ctx context.Context) error {
	select {
	case s.ch <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Release frees one slot. It panics if called more times than Acquire.
func (s *Semaphore) Release() {
	select {
	case <-s.ch:
	default:
		panic("semaphore: Release called without matching Acquire")
	}
}

// Cap returns the maximum concurrency allowed by this semaphore.
func (s *Semaphore) Cap() int {
	return cap(s.ch)
}

// Available returns the number of slots currently free.
func (s *Semaphore) Available() int {
	return cap(s.ch) - len(s.ch)
}
