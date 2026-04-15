package sync

import (
	"testing"
	"time"
)

func TestDefaultFenceConfig_Values(t *testing.T) {
	cfg := DefaultFenceConfig()
	if cfg.Window != 30*time.Second {
		t.Fatalf("expected 30s, got %v", cfg.Window)
	}
}

func TestNewWriteFence_ZeroWindowUsesDefault(t *testing.T) {
	f := NewWriteFence(FenceConfig{})
	if f.cfg.Window != 30*time.Second {
		t.Fatalf("expected default window, got %v", f.cfg.Window)
	}
}

func TestWriteFence_FirstCallAllowed(t *testing.T) {
	f := NewWriteFence(DefaultFenceConfig())
	if err := f.Allow("MY_KEY"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestWriteFence_SecondCallBlocked(t *testing.T) {
	f := NewWriteFence(DefaultFenceConfig())
	_ = f.Allow("MY_KEY")
	if err := f.Allow("MY_KEY"); err == nil {
		t.Fatal("expected fence error, got nil")
	}
}

func TestWriteFence_AllowsAfterWindowExpires(t *testing.T) {
	f := NewWriteFence(FenceConfig{Window: 10 * time.Millisecond})
	now := time.Now()
	f.nowFn = func() time.Time { return now }
	_ = f.Allow("MY_KEY")
	f.nowFn = func() time.Time { return now.Add(20 * time.Millisecond) }
	if err := f.Allow("MY_KEY"); err != nil {
		t.Fatalf("expected allow after window, got %v", err)
	}
}

func TestWriteFence_Reset_ClearsKey(t *testing.T) {
	f := NewWriteFence(DefaultFenceConfig())
	_ = f.Allow("MY_KEY")
	f.Reset("MY_KEY")
	if err := f.Allow("MY_KEY"); err != nil {
		t.Fatalf("expected allow after reset, got %v", err)
	}
}

func TestWriteFence_Flush_ClearsAll(t *testing.T) {
	f := NewWriteFence(DefaultFenceConfig())
	_ = f.Allow("A")
	_ = f.Allow("B")
	f.Flush()
	if err := f.Allow("A"); err != nil {
		t.Fatalf("expected allow after flush, got %v", err)
	}
	if err := f.Allow("B"); err != nil {
		t.Fatalf("expected allow after flush, got %v", err)
	}
}

func TestWriteFence_DifferentKeysIndependent(t *testing.T) {
	f := NewWriteFence(DefaultFenceConfig())
	_ = f.Allow("KEY_A")
	if err := f.Allow("KEY_B"); err != nil {
		t.Fatalf("KEY_B should not be fenced: %v", err)
	}
}
