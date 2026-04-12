package sync

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

func TestDefaultConcurrencyConfig_Values(t *testing.T) {
	cfg := DefaultConcurrencyConfig()
	if cfg.MaxWorkers != 4 {
		t.Errorf("expected MaxWorkers=4, got %d", cfg.MaxWorkers)
	}
	if cfg.FailFast {
		t.Error("expected FailFast=false")
	}
}

func TestRunConcurrent_AllSuccess(t *testing.T) {
	profiles := []string{"dev", "staging", "prod"}
	cfg := DefaultConcurrencyConfig()

	results := RunConcurrent(context.Background(), profiles, cfg, func(_ context.Context, _ string) error {
		return nil
	})

	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
	for _, r := range results {
		if r.Err != nil {
			t.Errorf("unexpected error for %s: %v", r.Profile, r.Err)
		}
	}
}

func TestRunConcurrent_PartialFailure(t *testing.T) {
	profiles := []string{"dev", "staging", "prod"}
	cfg := DefaultConcurrencyConfig()
	errBoom := errors.New("boom")

	results := RunConcurrent(context.Background(), profiles, cfg, func(_ context.Context, p string) error {
		if p == "staging" {
			return errBoom
		}
		return nil
	})

	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
	failed := 0
	for _, r := range results {
		if r.Err != nil {
			failed++
		}
	}
	if failed != 1 {
		t.Errorf("expected 1 failure, got %d", failed)
	}
}

func TestRunConcurrent_FailFast_CancelsOthers(t *testing.T) {
	profiles := []string{"a", "b", "c", "d", "e"}
	cfg := ConcurrencyConfig{MaxWorkers: 2, FailFast: true}
	errBoom := errors.New("boom")
	var executions int64

	RunConcurrent(context.Background(), profiles, cfg, func(ctx context.Context, p string) error {
		atomic.AddInt64(&executions, 1)
		if p == "a" {
			return errBoom
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(200 * time.Millisecond):
			return nil
		}
	})

	if atomic.LoadInt64(&executions) == int64(len(profiles)) {
		t.Error("expected FailFast to prevent all profiles from completing")
	}
}

func TestRunConcurrent_ZeroWorkers_DefaultsToOne(t *testing.T) {
	profiles := []string{"dev"}
	cfg := ConcurrencyConfig{MaxWorkers: 0}

	results := RunConcurrent(context.Background(), profiles, cfg, func(_ context.Context, _ string) error {
		return nil
	})

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
}

func TestRunConcurrent_EmptyProfiles(t *testing.T) {
	cfg := DefaultConcurrencyConfig()
	results := RunConcurrent(context.Background(), []string{}, cfg, func(_ context.Context, _ string) error {
		return nil
	})
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}
