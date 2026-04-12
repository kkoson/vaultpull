package sync

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestWithCircuitBreaker_AllowsWhenClosed(t *testing.T) {
	cb := NewCircuitBreaker(3, time.Second)
	called := false

	stage := WithCircuitBreaker(cb, func(_ context.Context, _ string) error {
		called = true
		return nil
	})

	if err := stage(context.Background(), "prod"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !called {
		t.Fatal("expected underlying stage to be called")
	}
}

func TestWithCircuitBreaker_BlocksWhenOpen(t *testing.T) {
	cb := NewCircuitBreaker(2, time.Second)

	// Trip the breaker.
	failStage := WithCircuitBreaker(cb, func(_ context.Context, _ string) error {
		return errors.New("vault unavailable")
	})
	_ = failStage(context.Background(), "prod")
	_ = failStage(context.Background(), "prod")

	called := false
	guardedStage := WithCircuitBreaker(cb, func(_ context.Context, _ string) error {
		called = true
		return nil
	})

	err := guardedStage(context.Background(), "prod")
	if err == nil {
		t.Fatal("expected error when circuit is open")
	}
	if called {
		t.Fatal("expected underlying stage NOT to be called when circuit is open")
	}
}

func TestWithCircuitBreaker_RecordsFailure(t *testing.T) {
	cb := NewCircuitBreaker(5, time.Second)
	expectedErr := errors.New("timeout")

	stage := WithCircuitBreaker(cb, func(_ context.Context, _ string) error {
		return expectedErr
	})

	err := stage(context.Background(), "staging")
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected %v, got %v", expectedErr, err)
	}
}

func TestWithCircuitBreaker_RecordsSuccess(t *testing.T) {
	cb := NewCircuitBreaker(3, time.Second)

	// Record a failure first.
	_ = WithCircuitBreaker(cb, func(_ context.Context, _ string) error {
		return errors.New("fail")
	})(context.Background(), "prod")

	// Then a success — breaker should remain closed.
	stage := WithCircuitBreaker(cb, func(_ context.Context, _ string) error {
		return nil
	})

	if err := stage(context.Background(), "prod"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should still be allowed (not open).
	if !cb.Allow() {
		t.Fatal("expected circuit to remain closed after success")
	}
}
