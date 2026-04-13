package sync

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestWithLease_ExecutesWhenNoLease(t *testing.T) {
	m := NewLeaseManager(DefaultLeaseConfig())
	called := false
	fn := WithLease(m, "dev", func(_ context.Context) error {
		called = true
		return nil
	})
	if err := fn(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Fatal("expected fn to be called")
	}
}

func TestWithLease_SkipsWhenLeaseValid(t *testing.T) {
	m := NewLeaseManager(DefaultLeaseConfig())
	// Acquire a fresh lease so renewal is not needed.
	m.Acquire("dev", time.Now())

	called := false
	fn := WithLease(m, "dev", func(_ context.Context) error {
		called = true
		return nil
	})
	if err := fn(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Fatal("expected fn to be skipped due to valid lease")
	}
}

func TestWithLease_ReleasesOnError(t *testing.T) {
	m := NewLeaseManager(DefaultLeaseConfig())
	expected := errors.New("vault unavailable")
	fn := WithLease(m, "prod", func(_ context.Context) error {
		return expected
	})
	err := fn(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, expected) {
		t.Fatalf("expected wrapped error, got %v", err)
	}
	// Lease should be released even on error.
	_, ok := m.Get("prod")
	if ok {
		t.Fatal("expected lease to be released after error")
	}
}

func TestWithLease_ExecutesWhenExpired(t *testing.T) {
	cfg := LeaseConfig{TTL: 1 * time.Millisecond, RenewThreshold: 0.99}
	m := NewLeaseManager(cfg)
	// Acquire then let it expire.
	past := time.Now().Add(-1 * time.Second)
	m.Acquire("staging", past)

	called := false
	fn := WithLease(m, "staging", func(_ context.Context) error {
		called = true
		return nil
	})
	if err := fn(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Fatal("expected fn to be called for expired lease")
	}
}
