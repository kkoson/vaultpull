package sync

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestNewLimiter_ZeroRate(t *testing.T) {
	l := NewLimiter(0)
	if l == nil {
		t.Fatal("expected non-nil limiter for zero rate")
	}
}

func TestNewLimiter_PositiveRate(t *testing.T) {
	l := NewLimiter(10)
	if l == nil {
		t.Fatal("expected non-nil limiter")
	}
}

func TestLimiter_Wait_ZeroRate_NoBlock(t *testing.T) {
	l := NewLimiter(0)
	ctx := context.Background()

	done := make(chan struct{})
	go func() {
		l.Wait(ctx)
		close(done)
	}()

	select {
	case <-done:
		// ok
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Wait blocked unexpectedly for zero-rate limiter")
	}
}

func TestLimiter_Wait_ContextCancelled(t *testing.T) {
	l := NewLimiter(1)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// Should return immediately when context is already cancelled
	done := make(chan struct{})
	go func() {
		l.Wait(ctx)
		close(done)
	}()

	select {
	case <-done:
		// ok
	case <-time.After(200 * time.Millisecond):
		t.Fatal("Wait did not respect cancelled context")
	}
}

func TestLimiter_Wait_RateLimit_Concurrent(t *testing.T) {
	const rps = 5
	l := NewLimiter(rps)
	ctx := context.Background()

	var wg sync.WaitGroup
	start := time.Now()

	for i := 0; i < rps; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			l.Wait(ctx)
		}()
	}

	wg.Wait()
	elapsed := time.Since(start)

	// All rps tokens should be consumed without excessive delay
	if elapsed > 2*time.Second {
		t.Errorf("rate limiter too slow: %v", elapsed)
	}
}
