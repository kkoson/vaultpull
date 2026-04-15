package sync

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

func TestDefaultBulkheadConfig_Values(t *testing.T) {
	cfg := DefaultBulkheadConfig()
	if cfg.MaxConcurrent != 10 {
		t.Errorf("expected MaxConcurrent=10, got %d", cfg.MaxConcurrent)
	}
	if cfg.MaxWaiting != 20 {
		t.Errorf("expected MaxWaiting=20, got %d", cfg.MaxWaiting)
	}
}

func TestNewBulkhead_ZeroFallsBackToDefaults(t *testing.T) {
	b := NewBulkhead(BulkheadConfig{})
	if b.cfg.MaxConcurrent != 10 {
		t.Errorf("expected MaxConcurrent=10, got %d", b.cfg.MaxConcurrent)
	}
}

func TestBulkhead_ExecutesSuccessfully(t *testing.T) {
	b := NewBulkhead(BulkheadConfig{MaxConcurrent: 2, MaxWaiting: 2})
	err := b.Execute(context.Background(), func(_ context.Context) error {
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBulkhead_PropagatesError(t *testing.T) {
	b := NewBulkhead(BulkheadConfig{MaxConcurrent: 1, MaxWaiting: 0})
	sentinel := errors.New("boom")
	err := b.Execute(context.Background(), func(_ context.Context) error {
		return sentinel
	})
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel error, got %v", err)
	}
}

func TestBulkhead_FullReturnsError(t *testing.T) {
	b := NewBulkhead(BulkheadConfig{MaxConcurrent: 1, MaxWaiting: 0})

	ready := make(chan struct{})
	done := make(chan struct{})
	go func() {
		_ = b.Execute(context.Background(), func(_ context.Context) error {
			close(ready)
			<-done
			return nil
		})
	}()
	<-ready

	err := b.Execute(context.Background(), func(_ context.Context) error { return nil })
	if !errors.Is(err, ErrBulkheadFull) {
		t.Fatalf("expected ErrBulkheadFull, got %v", err)
	}
	close(done)
}

func TestBulkhead_ContextCancelledWhileWaiting(t *testing.T) {
	b := NewBulkhead(BulkheadConfig{MaxConcurrent: 1, MaxWaiting: 5})

	blocking := make(chan struct{})
	go func() {
		_ = b.Execute(context.Background(), func(_ context.Context) error {
			<-blocking
			return nil
		})
	}()
	time.Sleep(10 * time.Millisecond)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := b.Execute(ctx, func(_ context.Context) error { return nil })
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
	close(blocking)
}

func TestBulkhead_Stats_ActiveAndWaiting(t *testing.T) {
	b := NewBulkhead(BulkheadConfig{MaxConcurrent: 1, MaxWaiting: 2})

	blocking := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		_ = b.Execute(context.Background(), func(_ context.Context) error {
			<-blocking
			return nil
		})
	}()
	time.Sleep(10 * time.Millisecond)
	active, _ := b.Stats()
	if active != 1 {
		t.Errorf("expected active=1, got %d", active)
	}
	close(blocking)
	wg.Wait()
}
