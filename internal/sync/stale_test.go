package sync

import (
	"testing"
	"time"
)

func TestDefaultStalenessConfig_Values(t *testing.T) {
	cfg := DefaultStalenessConfig()
	if cfg.MaxAge != 30*time.Minute {
		t.Errorf("expected 30m, got %v", cfg.MaxAge)
	}
	if cfg.WarnOnly {
		t.Error("expected WarnOnly=false")
	}
}

func TestNewStalenessTracker_ZeroMaxAgeUsesDefault(t *testing.T) {
	tracker := NewStalenessTracker(StalenessConfig{})
	if tracker.cfg.MaxAge != 30*time.Minute {
		t.Errorf("expected default MaxAge, got %v", tracker.cfg.MaxAge)
	}
}

func TestStalenessTracker_NeverSynced_IsStale(t *testing.T) {
	tracker := NewStalenessTracker(DefaultStalenessConfig())
	if !tracker.IsStale("prod") {
		t.Error("expected profile with no record to be stale")
	}
}

func TestStalenessTracker_RecentRecord_NotStale(t *testing.T) {
	tracker := NewStalenessTracker(StalenessConfig{MaxAge: 5 * time.Minute})
	tracker.Record("prod")
	if tracker.IsStale("prod") {
		t.Error("expected recently recorded profile to not be stale")
	}
}

func TestStalenessTracker_OldRecord_IsStale(t *testing.T) {
	tracker := NewStalenessTracker(StalenessConfig{MaxAge: time.Millisecond})
	tracker.Record("prod")
	time.Sleep(5 * time.Millisecond)
	if !tracker.IsStale("prod") {
		t.Error("expected old record to be stale")
	}
}

func TestStalenessTracker_LastRun_Missing(t *testing.T) {
	tracker := NewStalenessTracker(DefaultStalenessConfig())
	_, ok := tracker.LastRun("missing")
	if ok {
		t.Error("expected false for missing profile")
	}
}

func TestStalenessTracker_LastRun_Present(t *testing.T) {
	tracker := NewStalenessTracker(DefaultStalenessConfig())
	before := time.Now()
	tracker.Record("staging")
	t2, ok := tracker.LastRun("staging")
	if !ok {
		t.Fatal("expected record to exist")
	}
	if t2.Before(before) {
		t.Error("recorded time should not be before Record() was called")
	}
}

func TestStalenessTracker_Reset_MakesStale(t *testing.T) {
	tracker := NewStalenessTracker(StalenessConfig{MaxAge: 5 * time.Minute})
	tracker.Record("dev")
	tracker.Reset("dev")
	if !tracker.IsStale("dev") {
		t.Error("expected reset profile to be stale")
	}
}
