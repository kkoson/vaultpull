package sync

import (
	"testing"
	"time"
)

func TestDefaultBackoffConfig_Values(t *testing.T) {
	cfg := DefaultBackoffConfig()
	if cfg.Strategy != BackoffExponential {
		t.Errorf("expected BackoffExponential, got %v", cfg.Strategy)
	}
	if cfg.BaseDelay != 500*time.Millisecond {
		t.Errorf("unexpected BaseDelay: %v", cfg.BaseDelay)
	}
	if cfg.MaxDelay != 30*time.Second {
		t.Errorf("unexpected MaxDelay: %v", cfg.MaxDelay)
	}
	if cfg.Multiplier != 2.0 {
		t.Errorf("unexpected Multiplier: %v", cfg.Multiplier)
	}
}

func TestDelay_Fixed(t *testing.T) {
	cfg := BackoffConfig{
		Strategy:  BackoffFixed,
		BaseDelay: 1 * time.Second,
		MaxDelay:  10 * time.Second,
	}
	for _, attempt := range []int{0, 1, 5} {
		got := cfg.Delay(attempt)
		if got != 1*time.Second {
			t.Errorf("attempt %d: expected 1s, got %v", attempt, got)
		}
	}
}

func TestDelay_Linear(t *testing.T) {
	cfg := BackoffConfig{
		Strategy:  BackoffLinear,
		BaseDelay: 1 * time.Second,
		MaxDelay:  10 * time.Second,
	}
	cases := []struct{ attempt int; want time.Duration }{
		{0, 1 * time.Second},
		{1, 2 * time.Second},
		{4, 5 * time.Second},
	}
	for _, tc := range cases {
		got := cfg.Delay(tc.attempt)
		if got != tc.want {
			t.Errorf("attempt %d: expected %v, got %v", tc.attempt, tc.want, got)
		}
	}
}

func TestDelay_Exponential(t *testing.T) {
	cfg := BackoffConfig{
		Strategy:   BackoffExponential,
		BaseDelay:  1 * time.Second,
		MaxDelay:   60 * time.Second,
		Multiplier: 2.0,
	}
	cases := []struct{ attempt int; want time.Duration }{
		{0, 1 * time.Second},
		{1, 2 * time.Second},
		{2, 4 * time.Second},
		{3, 8 * time.Second},
	}
	for _, tc := range cases {
		got := cfg.Delay(tc.attempt)
		if got != tc.want {
			t.Errorf("attempt %d: expected %v, got %v", tc.attempt, tc.want, got)
		}
	}
}

func TestDelay_RespectsMaxDelay(t *testing.T) {
	cfg := BackoffConfig{
		Strategy:   BackoffExponential,
		BaseDelay:  1 * time.Second,
		MaxDelay:   5 * time.Second,
		Multiplier: 2.0,
	}
	got := cfg.Delay(10)
	if got != 5*time.Second {
		t.Errorf("expected max delay 5s, got %v", got)
	}
}

func TestDelay_NegativeAttemptClamped(t *testing.T) {
	cfg := BackoffConfig{
		Strategy:  BackoffFixed,
		BaseDelay: 2 * time.Second,
	}
	got := cfg.Delay(-3)
	if got != 2*time.Second {
		t.Errorf("expected 2s for negative attempt, got %v", got)
	}
}

func TestDelay_ZeroMultiplierDefaultsToTwo(t *testing.T) {
	cfg := BackoffConfig{
		Strategy:   BackoffExponential,
		BaseDelay:  1 * time.Second,
		MaxDelay:   60 * time.Second,
		Multiplier: 0,
	}
	got := cfg.Delay(3)
	if got != 8*time.Second {
		t.Errorf("expected 8s with default multiplier, got %v", got)
	}
}
