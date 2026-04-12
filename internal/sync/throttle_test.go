package sync

import (
	"context"
	"testing"
	"time"
)

func TestDefaultThrottleConfig_Values(t *testing.T) {
	cfg := DefaultThrottleConfig()
	if cfg.MinInterval != 200*time.Millisecond {
		t.Errorf("expected 200ms, got %v", cfg.MinInterval)
	}
	if cfg.BurstSize != 3 {
		t.Errorf("expected burst 3, got %d", cfg.BurstSize)
	}
}

func TestNewThrottle_ZeroIntervalUsesDefault(t *testing.T) {
	th := NewThrottle(ThrottleConfig{MinInterval: 0, BurstSize: 1})
	if th.cfg.MinInterval != DefaultThrottleConfig().MinInterval {
		t.Errorf("expected default interval, got %v", th.cfg.MinInterval)
	}
}

func TestNewThrottle_ZeroBurstDefaultsToOne(t *testing.T) {
	th := NewThrottle(ThrottleConfig{MinInterval: time.Millisecond, BurstSize: 0})
	if th.cfg.BurstSize != 1 {
		t.Errorf("expected burst 1, got %d", th.cfg.BurstSize)
	}
}

func TestThrottle_FirstCallDoesNotBlock(t *testing.T) {
	th := NewThrottle(ThrottleConfig{MinInterval: 10 * time.Second, BurstSize: 1})
	ctx := context.Background()
	start := time.Now()
	if err := th.Wait(ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if time.Since(start) > 50*time.Millisecond {
		t.Error("first call should not block")
	}
}

func TestThrottle_BurstAllowsMultipleCalls(t *testing.T) {
	th := NewThrottle(ThrottleConfig{MinInterval: 10 * time.Second, BurstSize: 3})
	ctx := context.Background()
	for i := 0; i < 3; i++ {
		if err := th.Wait(ctx); err != nil {
			t.Fatalf("call %d: unexpected error: %v", i, err)
		}
	}
}

func TestThrottle_ContextCancelledDuringWait(t *testing.T) {
	th := NewThrottle(ThrottleConfig{MinInterval: 5 * time.Second, BurstSize: 1})
	ctx, cancel := context.WithCancel(context.Background())
	// exhaust burst
	_ = th.Wait(ctx)
	cancel()
	err := th.Wait(ctx)
	if err == nil {
		t.Error("expected context error, got nil")
	}
}

func TestThrottle_Reset_ClearsState(t *testing.T) {
	th := NewThrottle(ThrottleConfig{MinInterval: 10 * time.Second, BurstSize: 1})
	ctx := context.Background()
	_ = th.Wait(ctx)
	th.Reset()
	start := time.Now()
	if err := th.Wait(ctx); err != nil {
		t.Fatalf("unexpected error after reset: %v", err)
	}
	if time.Since(start) > 50*time.Millisecond {
		t.Error("call after reset should not block")
	}
}
