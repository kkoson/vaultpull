package sync

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

func TestNewWorkerPool_DefaultsToOne(t *testing.T) {
	p := NewWorkerPool(0)
	if p.workers != 1 {
		t.Fatalf("expected 1 worker, got %d", p.workers)
	}
}

func TestNewWorkerPool_ValidLimit(t *testing.T) {
	p := NewWorkerPool(4)
	if p.workers != 4 {
		t.Fatalf("expected 4 workers, got %d", p.workers)
	}
}

func TestWorkerPool_AllJobsCompleted(t *testing.T) {
	ctx := context.Background()
	p := NewWorkerPool(3)
	results := p.Start(ctx)

	profiles := []string{"dev", "staging", "prod"}
	for _, name := range profiles {
		name := name
		if err := p.Submit(ctx, name, func(ctx context.Context) error {
			return nil
		}); err != nil {
			t.Fatalf("submit %q: %v", name, err)
		}
	}
	p.Close()

	var got []string
	for r := range results {
		if r.Err != nil {
			t.Errorf("unexpected error for %q: %v", r.ProfileName, r.Err)
		}
		got = append(got, r.ProfileName)
	}
	if len(got) != len(profiles) {
		t.Fatalf("expected %d results, got %d", len(profiles), len(got))
	}
}

func TestWorkerPool_ErrorPropagated(t *testing.T) {
	ctx := context.Background()
	p := NewWorkerPool(2)
	results := p.Start(ctx)

	wantErr := errors.New("vault unreachable")
	_ = p.Submit(ctx, "prod", func(ctx context.Context) error { return wantErr })
	p.Close()

	for r := range results {
		if r.ProfileName == "prod" && !errors.Is(r.Err, wantErr) {
			t.Errorf("expected %v, got %v", wantErr, r.Err)
		}
	}
}

func TestWorkerPool_ContextCancelledSubmit(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	p := NewWorkerPool(1)
	// Do not call Start so jobs channel is never drained.
	err := p.Submit(ctx, "dev", func(ctx context.Context) error { return nil })
	if err == nil {
		t.Fatal("expected error on cancelled context submit")
	}
}

func TestWorkerPool_ConcurrentExecution(t *testing.T) {
	ctx := context.Background()
	p := NewWorkerPool(4)
	results := p.Start(ctx)

	var concurrent int64
	var peak int64
	for i := 0; i < 8; i++ {
		_ = p.Submit(ctx, "p", func(ctx context.Context) error {
			cur := atomic.AddInt64(&concurrent, 1)
			for {
				old := atomic.LoadInt64(&peak)
				if cur <= old || atomic.CompareAndSwapInt64(&peak, old, cur) {
					break
				}
			}
			time.Sleep(5 * time.Millisecond)
			atomic.AddInt64(&concurrent, -1)
			return nil
		})
	}
	p.Close()
	for range results {
	}
	if peak < 2 {
		t.Errorf("expected concurrent execution, peak was %d", peak)
	}
}
