package sync

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/yourusername/vaultpull/internal/config"
)

func buildExpiryProfile(name string) config.Profile {
	return config.Profile{Name: name, VaultPath: "secret/data/" + name, OutputFile: name + ".env"}
}

func TestWithExpiry_PanicsOnNil(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic on nil tracker")
		}
	}()
	WithExpiry(nil, func(_ context.Context, _ config.Profile) error { return nil })
}

func TestWithExpiry_AllowsWhenExpired(t *testing.T) {
	tracker := NewExpiryTracker(ExpiryConfig{TTL: time.Millisecond, CleanupPeriod: time.Minute})
	called := false
	fn := WithExpiry(tracker, func(_ context.Context, _ config.Profile) error {
		called = true
		return nil
	})
	// Never recorded → expired → should call through.
	if err := fn(context.Background(), buildExpiryProfile("prod")); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !called {
		t.Error("expected inner func to be called")
	}
}

func TestWithExpiry_SkipsWhenFresh(t *testing.T) {
	tracker := NewExpiryTracker(ExpiryConfig{TTL: time.Hour, CleanupPeriod: time.Minute})
	p := buildExpiryProfile("staging")
	tracker.Record(p.Name)

	called := false
	fn := WithExpiry(tracker, func(_ context.Context, _ config.Profile) error {
		called = true
		return nil
	})
	err := fn(context.Background(), p)
	if err == nil || !strings.Contains(err.Error(), "still fresh") {
		t.Errorf("expected fresh-skip error, got %v", err)
	}
	if called {
		t.Error("expected inner func NOT to be called")
	}
}

func TestWithExpiry_RecordsOnSuccess(t *testing.T) {
	tracker := NewExpiryTracker(ExpiryConfig{TTL: time.Hour, CleanupPeriod: time.Minute})
	p := buildExpiryProfile("dev")

	fn := WithExpiry(tracker, func(_ context.Context, _ config.Profile) error { return nil })
	if err := fn(context.Background(), p); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tracker.IsExpired(p.Name) {
		t.Error("expected profile to be recorded as fresh after success")
	}
}

func TestWithExpiry_PropagatesInnerError(t *testing.T) {
	tracker := NewExpiryTracker(ExpiryConfig{TTL: time.Millisecond, CleanupPeriod: time.Minute})
	sentinel := errors.New("vault down")
	p := buildExpiryProfile("prod")

	fn := WithExpiry(tracker, func(_ context.Context, _ config.Profile) error { return sentinel })
	err := fn(context.Background(), p)
	if !errors.Is(err, sentinel) {
		t.Errorf("expected sentinel error, got %v", err)
	}
	// Should NOT be recorded on failure.
	time.Sleep(2 * time.Millisecond)
	if !tracker.IsExpired(p.Name) {
		t.Error("expected profile to remain expired after inner error")
	}
}
