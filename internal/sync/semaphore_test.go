package sync

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestNewSemaphore_InvalidLimit(t *testing.T) {
	_, err := NewSemaphore(0)
	if err == nil {
		t.Fatal("expected error for zero limit")
	}
	_, err = NewSemaphore(-3)
	if err == nil {
		t.Fatal("expected error for negative limit")
	}
}

func TestNewSemaphore_ValidLimit(t *testing.T) {
	s, err := NewSemaphore(3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Cap() != 3 {
		t.Errorf("expected cap 3, got %d", s.Cap())
	}
	if s.Available() != 3 {
		t.Errorf("expected 3 available slots, got %d", s.Available())
	}
}

func TestSemaphore_AcquireRelease(t *testing.T) {
	s, _ := NewSemaphore(2)
	ctx := context.Background()

	if err := s.Acquire(ctx); err != nil {
		t.Fatalf("first acquire failed: %v", err)
	}
	if s.Available() != 1 {
		t.Errorf("expected 1 available, got %d", s.Available())
	}
	s.Release()
	if s.Available() != 2 {
		t.Errorf("expected 2 available after release, got %d", s.Available())
	}
}

func TestSemaphore_BlocksAtLimit(t *testing.T) {
	s, _ := NewSemaphore(1)
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	// Fill the single slot.
	_ = s.Acquire(context.Background())

	// Second acquire should block and then return ctx error.
	err := s.Acquire(ctx)
	if err == nil {
		t.Fatal("expected context error when semaphore is full")
	}
}

func TestSemaphore_ConcurrentAcquires(t *testing.T) {
	const workers = 10
	const limit = 3
	s, _ := NewSemaphore(limit)
	ctx := context.Background()

	var wg sync.WaitGroup
	var mu sync.Mutex
	maxConcurrent := 0
	active := 0

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = s.Acquire(ctx)
			mu.Lock()
			active++
			if active > maxConcurrent {
				maxConcurrent = active
			}
			mu.Unlock()
			time.Sleep(5 * time.Millisecond)
			mu.Lock()
			active--
			mu.Unlock()
			s.Release()
		}()
	}
	wg.Wait()

	if maxConcurrent > limit {
		t.Errorf("concurrency exceeded limit: max=%d limit=%d", maxConcurrent, limit)
	}
}

func TestSemaphore_ReleaseWithoutAcquirePanics(t *testing.T) {
	s, _ := NewSemaphore(2)
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic on Release without Acquire")
		}
	}()
	s.Release()
}
