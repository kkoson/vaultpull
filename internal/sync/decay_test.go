package sync

import (
	"testing"
	"time"
)

func TestDefaultDecayConfig_Values(t *testing.T) {
	cfg := DefaultDecayConfig()
	if cfg.HalfLife != 5*time.Minute {
		t.Errorf("expected 5m half-life, got %v", cfg.HalfLife)
	}
}

func TestNewDecayCounter_ZeroHalfLifeUsesDefault(t *testing.T) {
	d := NewDecayCounter(DecayConfig{})
	if d.cfg.HalfLife != DefaultDecayConfig().HalfLife {
		t.Error("expected default half-life to be applied")
	}
}

func TestDecayCounter_InitialValueIsZero(t *testing.T) {
	d := NewDecayCounter(DefaultDecayConfig())
	if d.Value() != 0 {
		t.Errorf("expected 0, got %f", d.Value())
	}
}

func TestDecayCounter_AddIncrements(t *testing.T) {
	d := NewDecayCounter(DefaultDecayConfig())
	d.Add(10)
	v := d.Value()
	// Value should be close to 10 immediately after adding.
	if v < 9.9 || v > 10.1 {
		t.Errorf("expected ~10, got %f", v)
	}
}

func TestDecayCounter_ValueDecaysOverTime(t *testing.T) {
	// Use a very short half-life so decay is observable in tests.
	cfg := DecayConfig{HalfLife: 50 * time.Millisecond}
	d := NewDecayCounter(cfg)
	d.Add(100)

	time.Sleep(100 * time.Millisecond) // two half-lives → ~25
	v := d.Value()
	if v >= 50 {
		t.Errorf("expected value to decay below 50, got %f", v)
	}
}

func TestDecayCounter_Reset_ZeroesValue(t *testing.T) {
	d := NewDecayCounter(DefaultDecayConfig())
	d.Add(42)
	d.Reset()
	if d.Value() != 0 {
		t.Errorf("expected 0 after reset, got %f", d.Value())
	}
}

func TestDecayCounter_MultipleAdds_Accumulate(t *testing.T) {
	cfg := DecayConfig{HalfLife: 10 * time.Second}
	d := NewDecayCounter(cfg)
	d.Add(5)
	d.Add(5)
	v := d.Value()
	// Both adds happen nearly simultaneously so value should be ~10.
	if v < 9.5 {
		t.Errorf("expected ~10, got %f", v)
	}
}
