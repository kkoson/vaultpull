package sync

import (
	"testing"
	"time"
)

func TestNewRateLimitStore_ZeroRateAlwaysAllows(t *testing.T) {
	s := NewRateLimitStore(0, time.Minute)
	for i := 0; i < 100; i++ {
		if !s.Allow("profile-a") {
			t.Fatal("expected allow with zero maxRate")
		}
	}
}

func TestRateLimitStore_AllowsUpToMaxRate(t *testing.T) {
	s := NewRateLimitStore(3, time.Minute)
	for i := 0; i < 3; i++ {
		if !s.Allow("p") {
			t.Fatalf("expected allow on attempt %d", i+1)
		}
	}
	if s.Allow("p") {
		t.Fatal("expected deny after maxRate exceeded")
	}
}

func TestRateLimitStore_DeniedCountIncremented(t *testing.T) {
	s := NewRateLimitStore(1, time.Minute)
	s.Allow("p") // allowed
	s.Allow("p") // denied
	s.Allow("p") // denied

	stats := s.Stats("p")
	if stats == nil {
		t.Fatal("expected stats")
	}
	if stats.Denied != 2 {
		t.Errorf("expected 2 denied, got %d", stats.Denied)
	}
	if stats.Allowed != 1 {
		t.Errorf("expected 1 allowed, got %d", stats.Allowed)
	}
}

func TestRateLimitStore_WindowReset(t *testing.T) {
	s := NewRateLimitStore(1, 10*time.Millisecond)
	if !s.Allow("p") {
		t.Fatal("first allow should succeed")
	}
	if s.Allow("p") {
		t.Fatal("second allow within window should fail")
	}
	time.Sleep(20 * time.Millisecond)
	if !s.Allow("p") {
		t.Fatal("allow after window reset should succeed")
	}
}

func TestRateLimitStore_Stats_MissingProfile(t *testing.T) {
	s := NewRateLimitStore(5, time.Minute)
	if s.Stats("nonexistent") != nil {
		t.Fatal("expected nil stats for unknown profile")
	}
}

func TestRateLimitStore_Reset_ClearsState(t *testing.T) {
	s := NewRateLimitStore(1, time.Minute)
	s.Allow("p")
	s.Allow("p") // denied
	s.Reset("p")
	if !s.Allow("p") {
		t.Fatal("expected allow after reset")
	}
}

func TestRateLimitStore_IndependentProfiles(t *testing.T) {
	s := NewRateLimitStore(1, time.Minute)
	if !s.Allow("a") {
		t.Fatal("profile a first allow")
	}
	if !s.Allow("b") {
		t.Fatal("profile b first allow")
	}
	if s.Allow("a") {
		t.Fatal("profile a should be denied")
	}
	if s.Allow("b") {
		t.Fatal("profile b should be denied")
	}
}
