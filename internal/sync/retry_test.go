package sync

import (
	"context"
	"errors"
	"testing"
	"time"
)

var errTemp = errors.New("temporary error")

func TestWithRetry_SuccessOnFirstAttempt(t *testing.T) {
	calls := 0
	err := WithRetry(context.Background(), DefaultRetryPolicy(), func() error {
		calls++
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestWithRetry_SuccessOnSecondAttempt(t *testing.T) {
	calls := 0
	policy := RetryPolicy{MaxAttempts: 3, Delay: time.Millisecond, Multiplier: 1.0}
	err := WithRetry(context.Background(), policy, func() error {
		calls++
		if calls < 2 {
			return errTemp
		}
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if calls != 2 {
		t.Fatalf("expected 2 calls, got %d", calls)
	}
}

func TestWithRetry_ExhaustsAttempts(t *testing.T) {
	calls := 0
	policy := RetryPolicy{MaxAttempts: 3, Delay: time.Millisecond, Multiplier: 1.0}
	err := WithRetry(context.Background(), policy, func() error {
		calls++
		return errTemp
	})
	if !errors.Is(err, ErrMaxAttemptsReached) {
		t.Fatalf("expected ErrMaxAttemptsReached, got %v", err)
	}
	if !errors.Is(err, errTemp) {
		t.Fatalf("expected wrapped errTemp, got %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestWithRetry_ContextCancelledBeforeStart(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	calls := 0
	err := WithRetry(ctx, DefaultRetryPolicy(), func() error {
		calls++
		return nil
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
	if calls != 0 {
		t.Fatalf("expected 0 calls, got %d", calls)
	}
}

func TestWithRetry_ContextCancelledDuringDelay(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	policy := RetryPolicy{MaxAttempts: 5, Delay: 200 * time.Millisecond, Multiplier: 1.0}

	calls := 0
	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	err := WithRetry(ctx, policy, func() error {
		calls++
		return errTemp
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}

func TestDefaultRetryPolicy_Values(t *testing.T) {
	p := DefaultRetryPolicy()
	if p.MaxAttempts != 3 {
		t.Errorf("expected MaxAttempts=3, got %d", p.MaxAttempts)
	}
	if p.Multiplier != 2.0 {
		t.Errorf("expected Multiplier=2.0, got %f", p.Multiplier)
	}
	if p.Delay != 500*time.Millisecond {
		t.Errorf("expected Delay=500ms, got %v", p.Delay)
	}
}
