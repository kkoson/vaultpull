package sync

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/yourusername/vaultpull/internal/config"
)

func TestWithBulkhead_PanicsOnNil(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for nil bulkhead")
		}
	}()
	WithBulkhead(nil)
}

func TestWithBulkhead_AllowsExecution(t *testing.T) {
	b := NewBulkhead(BulkheadConfig{MaxConcurrent: 2, MaxWaiting: 2})
	p := config.Profile{Name: "dev"}

	var called bool
	stage := WithBulkhead(b)(func(_ context.Context, _ config.Profile) error {
		called = true
		return nil
	})

	if err := stage(context.Background(), p); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Error("expected inner stage to be called")
	}
}

func TestWithBulkhead_WrapsError(t *testing.T) {
	b := NewBulkhead(BulkheadConfig{MaxConcurrent: 1, MaxWaiting: 0})
	p := config.Profile{Name: "staging"}

	blocking := make(chan struct{})
	go func() {
		_ = b.Execute(context.Background(), func(_ context.Context) error {
			<-blocking
			return nil
		})
	}()
	time.Sleep(10 * time.Millisecond)

	sentinel := errors.New("inner error")
	stage := WithBulkhead(b)(func(_ context.Context, _ config.Profile) error {
		return sentinel
	})

	err := stage(context.Background(), p)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, ErrBulkheadFull) {
		t.Errorf("expected ErrBulkheadFull in error chain, got %v", err)
	}
	close(blocking)
}

func TestWithBulkhead_PropagatesInnerError(t *testing.T) {
	b := NewBulkhead(BulkheadConfig{MaxConcurrent: 2, MaxWaiting: 2})
	p := config.Profile{Name: "prod"}
	sentinel := errors.New("write failed")

	stage := WithBulkhead(b)(func(_ context.Context, _ config.Profile) error {
		return sentinel
	})

	err := stage(context.Background(), p)
	if !errors.Is(err, sentinel) {
		t.Errorf("expected sentinel in error chain, got %v", err)
	}
}
