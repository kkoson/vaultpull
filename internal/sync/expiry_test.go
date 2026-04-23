package sync

import (
	"testing"
	"time"
)

func TestDefaultExpiryConfig_Values(t *testing.T) {
	cfg := DefaultExpiryConfig()
	if cfg.TTL != 30*time.Minute {
		t.Errorf("expected TTL 30m, got %v", cfg.TTL)
	}
	if cfg.CleanupPeriod != 5*time.Minute {
		t.Errorf("expected CleanupPeriod 5m, got %v", cfg.CleanupPeriod)
	}
}

func TestNewExpiryTracker_ZeroConfigUsesDefaults(t *testing.T) {
	tr := NewExpiryTracker(ExpiryConfig{})
	if tr.cfg.TTL != 30*time.Minute {
		t.Errorf("expected default TTL, got %v", tr.cfg.TTL)
	}
}

func TestExpiryTracker_NeverRecorded_IsExpired(t *testing.T) {
	tr := NewExpiryTracker(DefaultExpiryConfig())
	if !tr.IsExpired("prod") {
		t.Error("expected unrecorded profile to be expired")
	}
}

func TestExpiryTracker_RecentRecord_NotExpired(t *testing.T) {
	tr := NewExpiryTracker(ExpiryConfig{TTL: time.Hour, CleanupPeriod: time.Minute})
	tr.Record("prod")
	if tr.IsExpired("prod") {
		t.Error("expected freshly recorded profile to not be expired")
	}
}

func TestExpiryTracker_OldRecord_IsExpired(t *testing.T) {
	tr := NewExpiryTracker(ExpiryConfig{TTL: time.Millisecond, CleanupPeriod: time.Minute})
	tr.Record("prod")
	time.Sleep(5 * time.Millisecond)
	if !tr.IsExpired("prod") {
		t.Error("expected old record to be expired")
	}
}

func TestExpiryTracker_Len_IncrementsOnRecord(t *testing.T) {
	tr := NewExpiryTracker(DefaultExpiryConfig())
	if tr.Len() != 0 {
		t.Errorf("expected 0, got %d", tr.Len())
	}
	tr.Record("a")
	tr.Record("b")
	if tr.Len() != 2 {
		t.Errorf("expected 2, got %d", tr.Len())
	}
}

func TestExpiryTracker_Evict_RemovesExpired(t *testing.T) {
	tr := NewExpiryTracker(ExpiryConfig{TTL: time.Millisecond, CleanupPeriod: time.Minute})
	tr.Record("old")
	time.Sleep(5 * time.Millisecond)
	tr.Record("fresh") // re-record to reset its timestamp
	// overwrite fresh with a long TTL tracker — use separate tracker
	tr2 := NewExpiryTracker(ExpiryConfig{TTL: time.Hour, CleanupPeriod: time.Minute})
	tr2.Record("fresh")

	removed := tr.Evict()
	if removed != 1 {
		t.Errorf("expected 1 evicted, got %d", removed)
	}
	if tr.Len() != 0 {
		t.Errorf("expected 0 remaining, got %d", tr.Len())
	}
}

func TestExpiryTracker_Evict_KeepsFreshEntries(t *testing.T) {
	tr := NewExpiryTracker(ExpiryConfig{TTL: time.Hour, CleanupPeriod: time.Minute})
	tr.Record("a")
	tr.Record("b")
	removed := tr.Evict()
	if removed != 0 {
		t.Errorf("expected 0 evicted, got %d", removed)
	}
	if tr.Len() != 2 {
		t.Errorf("expected 2 remaining, got %d", tr.Len())
	}
}
