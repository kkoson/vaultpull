package sync

import (
	"testing"
	"time"
)

func TestDefaultCooldownConfig_Values(t *testing.T) {
	cfg := DefaultCooldownConfig()
	if cfg.Duration != 30*time.Second {
		t.Errorf("expected 30s, got %v", cfg.Duration)
	}
}

func TestNewCooldownManager_ZeroDurationUsesDefault(t *testing.T) {
	cm := NewCooldownManager(CooldownConfig{Duration: 0})
	if cm.cfg.Duration != 30*time.Second {
		t.Errorf("expected default 30s, got %v", cm.cfg.Duration)
	}
}

func TestCooldownManager_FirstCallAllowed(t *testing.T) {
	cm := NewCooldownManager(CooldownConfig{Duration: 5 * time.Second})
	if !cm.Allow("prod") {
		t.Error("expected first call to be allowed")
	}
}

func TestCooldownManager_SecondCallBlocked(t *testing.T) {
	cm := NewCooldownManager(CooldownConfig{Duration: 5 * time.Second})
	cm.Allow("prod")
	if cm.Allow("prod") {
		t.Error("expected second immediate call to be blocked")
	}
}

func TestCooldownManager_AllowsAfterCooldown(t *testing.T) {
	cm := NewCooldownManager(CooldownConfig{Duration: 10 * time.Millisecond})
	cm.Allow("staging")
	time.Sleep(20 * time.Millisecond)
	if !cm.Allow("staging") {
		t.Error("expected call to be allowed after cooldown expired")
	}
}

func TestCooldownManager_Reset_AllowsImmediately(t *testing.T) {
	cm := NewCooldownManager(CooldownConfig{Duration: 5 * time.Second})
	cm.Allow("dev")
	cm.Reset("dev")
	if !cm.Allow("dev") {
		t.Error("expected call to be allowed after reset")
	}
}

func TestCooldownManager_Remaining_ZeroWhenNotTracked(t *testing.T) {
	cm := NewCooldownManager(CooldownConfig{Duration: 5 * time.Second})
	if r := cm.Remaining("unknown"); r != 0 {
		t.Errorf("expected 0 remaining, got %v", r)
	}
}

func TestCooldownManager_Remaining_PositiveInWindow(t *testing.T) {
	cm := NewCooldownManager(CooldownConfig{Duration: 5 * time.Second})
	cm.Allow("prod")
	r := cm.Remaining("prod")
	if r <= 0 || r > 5*time.Second {
		t.Errorf("expected remaining between 0 and 5s, got %v", r)
	}
}

func TestCooldownManager_IndependentProfiles(t *testing.T) {
	cm := NewCooldownManager(CooldownConfig{Duration: 5 * time.Second})
	cm.Allow("prod")
	if !cm.Allow("staging") {
		t.Error("expected staging to be allowed independently of prod")
	}
}
