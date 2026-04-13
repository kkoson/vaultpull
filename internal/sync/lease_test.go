package sync

import (
	"testing"
	"time"
)

func TestDefaultLeaseConfig_Values(t *testing.T) {
	cfg := DefaultLeaseConfig()
	if cfg.TTL != 5*time.Minute {
		t.Fatalf("expected 5m TTL, got %v", cfg.TTL)
	}
	if cfg.RenewThreshold != 0.25 {
		t.Fatalf("expected 0.25 threshold, got %v", cfg.RenewThreshold)
	}
}

func TestLease_IsExpired(t *testing.T) {
	now := time.Now()
	l := Lease{AcquiredAt: now.Add(-2 * time.Minute), ExpiresAt: now.Add(-1 * time.Second)}
	if !l.IsExpired(now) {
		t.Fatal("expected lease to be expired")
	}
	l2 := Lease{AcquiredAt: now, ExpiresAt: now.Add(time.Minute)}
	if l2.IsExpired(now) {
		t.Fatal("expected lease to be valid")
	}
}

func TestLease_NeedsRenewal_BelowThreshold(t *testing.T) {
	now := time.Now()
	// 10% remaining of 100s => needs renewal at 25% threshold
	l := Lease{
		AcquiredAt: now.Add(-90 * time.Second),
		ExpiresAt:  now.Add(10 * time.Second),
	}
	if !l.NeedsRenewal(now, 0.25) {
		t.Fatal("expected renewal needed")
	}
}

func TestLease_NeedsRenewal_AboveThreshold(t *testing.T) {
	now := time.Now()
	// 80% remaining => no renewal at 25% threshold
	l := Lease{
		AcquiredAt: now.Add(-20 * time.Second),
		ExpiresAt:  now.Add(80 * time.Second),
	}
	if l.NeedsRenewal(now, 0.25) {
		t.Fatal("expected no renewal needed")
	}
}

func TestLeaseManager_AcquireAndGet(t *testing.T) {
	m := NewLeaseManager(DefaultLeaseConfig())
	now := time.Now()
	l := m.Acquire("dev", now)
	if l.Profile != "dev" {
		t.Fatalf("expected profile dev, got %s", l.Profile)
	}
	got, ok := m.Get("dev")
	if !ok {
		t.Fatal("expected lease to exist")
	}
	if got.ExpiresAt != l.ExpiresAt {
		t.Fatal("stored lease mismatch")
	}
}

func TestLeaseManager_Release(t *testing.T) {
	m := NewLeaseManager(DefaultLeaseConfig())
	m.Acquire("prod", time.Now())
	m.Release("prod")
	_, ok := m.Get("prod")
	if ok {
		t.Fatal("expected lease to be removed")
	}
}

func TestLeaseManager_NeedsRenewal_NoLease(t *testing.T) {
	m := NewLeaseManager(DefaultLeaseConfig())
	if !m.NeedsRenewal("missing", time.Now()) {
		t.Fatal("expected renewal needed when no lease exists")
	}
}

func TestLeaseManager_NeedsRenewal_FreshLease(t *testing.T) {
	m := NewLeaseManager(DefaultLeaseConfig())
	now := time.Now()
	m.Acquire("dev", now)
	if m.NeedsRenewal("dev", now) {
		t.Fatal("expected no renewal needed for fresh lease")
	}
}

func TestNewLeaseManager_DefaultsOnZero(t *testing.T) {
	m := NewLeaseManager(LeaseConfig{})
	if m.cfg.TTL != 5*time.Minute {
		t.Fatalf("expected default TTL, got %v", m.cfg.TTL)
	}
	if m.cfg.RenewThreshold != 0.25 {
		t.Fatalf("expected default threshold, got %v", m.cfg.RenewThreshold)
	}
}
