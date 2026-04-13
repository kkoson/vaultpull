package sync

import (
	"testing"
	"time"
)

func TestDefaultJitterConfig_Values(t *testing.T) {
	cfg := DefaultJitterConfig()
	if cfg.Factor != 0.2 {
		t.Errorf("expected Factor 0.2, got %v", cfg.Factor)
	}
	if cfg.MaxJitter != 5*time.Second {
		t.Errorf("expected MaxJitter 5s, got %v", cfg.MaxJitter)
	}
}

func TestJitter_ZeroBase_ReturnsZero(t *testing.T) {
	result := Jitter(0, nil)
	if result != 0 {
		t.Errorf("expected 0, got %v", result)
	}
}

func TestJitter_NegativeBase_ReturnsZero(t *testing.T) {
	result := Jitter(-1*time.Second, nil)
	if result != 0 {
		t.Errorf("expected 0, got %v", result)
	}
}

func TestJitter_ZeroFactor_ReturnsBase(t *testing.T) {
	cfg := &JitterConfig{Factor: 0, MaxJitter: 0}
	base := 2 * time.Second
	result := Jitter(base, cfg)
	if result != base {
		t.Errorf("expected %v, got %v", base, result)
	}
}

func TestJitter_NilConfig_UsesDefaults(t *testing.T) {
	base := 1 * time.Second
	for i := 0; i < 50; i++ {
		result := Jitter(base, nil)
		// With factor 0.2 and MaxJitter 5s, result should be within ±10% of base
		low := time.Duration(float64(base) * 0.9)
		high := time.Duration(float64(base) * 1.1)
		if result < low || result > high {
			t.Errorf("iteration %d: jitter %v out of expected range [%v, %v]", i, result, low, high)
		}
	}
}

func TestJitter_RespectsMaxJitter(t *testing.T) {
	cfg := &JitterConfig{
		Factor:    1.0,
		MaxJitter: 100 * time.Millisecond,
	}
	base := 10 * time.Second
	for i := 0; i < 50; i++ {
		result := Jitter(base, cfg)
		low := base - 100*time.Millisecond
		high := base + 100*time.Millisecond
		if result < low || result > high {
			t.Errorf("iteration %d: result %v outside capped range [%v, %v]", i, result, low, high)
		}
	}
}

func TestJitter_InvalidFactor_FallsBackToDefault(t *testing.T) {
	cfg := &JitterConfig{Factor: 1.5} // invalid, > 1
	base := 1 * time.Second
	// Should not panic; factor clamped to default 0.2
	result := Jitter(base, cfg)
	if result <= 0 {
		t.Errorf("expected positive duration, got %v", result)
	}
}
