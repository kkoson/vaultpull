package sync

import (
	"testing"
	"time"
)

func TestDefaultQuotaConfig_Values(t *testing.T) {
	cfg := DefaultQuotaConfig()
	if cfg.MaxSyncsPerWindow != 100 {
		t.Errorf("expected 100, got %d", cfg.MaxSyncsPerWindow)
	}
	if cfg.Window != time.Hour {
		t.Errorf("expected 1h, got %s", cfg.Window)
	}
}

func TestNewQuotaEnforcer_DefaultsOnZero(t *testing.T) {
	q := NewQuotaEnforcer(QuotaConfig{})
	if q.cfg.MaxSyncsPerWindow != 100 {
		t.Errorf("expected default max 100, got %d", q.cfg.MaxSyncsPerWindow)
	}
	if q.cfg.Window != time.Hour {
		t.Errorf("expected default window 1h, got %s", q.cfg.Window)
	}
}

func TestQuotaEnforcer_AllowsWithinLimit(t *testing.T) {
	q := NewQuotaEnforcer(QuotaConfig{MaxSyncsPerWindow: 3, Window: time.Minute})
	for i := 0; i < 3; i++ {
		if err := q.Allow("dev"); err != nil {
			t.Fatalf("expected allow on call %d, got error: %v", i+1, err)
		}
	}
}

func TestQuotaEnforcer_BlocksOverLimit(t *testing.T) {
	q := NewQuotaEnforcer(QuotaConfig{MaxSyncsPerWindow: 2, Window: time.Minute})
	_ = q.Allow("dev")
	_ = q.Allow("dev")
	if err := q.Allow("dev"); err == nil {
		t.Fatal("expected quota exceeded error, got nil")
	}
}

func TestQuotaEnforcer_WindowReset(t *testing.T) {
	now := time.Now()
	q := NewQuotaEnforcer(QuotaConfig{MaxSyncsPerWindow: 1, Window: time.Millisecond})
	q.now = func() time.Time { return now }
	_ = q.Allow("dev")

	// Advance time past the window.
	q.now = func() time.Time { return now.Add(2 * time.Millisecond) }
	if err := q.Allow("dev"); err != nil {
		t.Fatalf("expected allow after window reset, got: %v", err)
	}
}

func TestQuotaEnforcer_Stats_NoUsage(t *testing.T) {
	q := NewQuotaEnforcer(DefaultQuotaConfig())
	count, end := q.Stats("prod")
	if count != 0 {
		t.Errorf("expected 0 count, got %d", count)
	}
	if !end.IsZero() {
		t.Errorf("expected zero time, got %s", end)
	}
}

func TestQuotaEnforcer_Stats_AfterAllow(t *testing.T) {
	q := NewQuotaEnforcer(QuotaConfig{MaxSyncsPerWindow: 10, Window: time.Hour})
	_ = q.Allow("staging")
	_ = q.Allow("staging")
	count, end := q.Stats("staging")
	if count != 2 {
		t.Errorf("expected count 2, got %d", count)
	}
	if end.IsZero() {
		t.Error("expected non-zero window end")
	}
}

func TestQuotaEnforcer_Reset_ClearsEntry(t *testing.T) {
	q := NewQuotaEnforcer(QuotaConfig{MaxSyncsPerWindow: 1, Window: time.Hour})
	_ = q.Allow("dev")
	q.Reset("dev")
	if err := q.Allow("dev"); err != nil {
		t.Fatalf("expected allow after reset, got: %v", err)
	}
}

func TestQuotaEnforcer_IsolatesProfiles(t *testing.T) {
	q := NewQuotaEnforcer(QuotaConfig{MaxSyncsPerWindow: 1, Window: time.Hour})
	_ = q.Allow("dev")
	if err := q.Allow("prod"); err != nil {
		t.Fatalf("prod should be independent of dev quota, got: %v", err)
	}
}
