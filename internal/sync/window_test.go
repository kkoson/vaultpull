package sync

import (
	"testing"
	"time"
)

func TestDefaultWindowConfig_Values(t *testing.T) {
	cfg := DefaultWindowConfig()
	if cfg.Size != time.Minute {
		t.Errorf("expected Size=1m, got %v", cfg.Size)
	}
	if cfg.BucketCount != 10 {
		t.Errorf("expected BucketCount=10, got %d", cfg.BucketCount)
	}
}

func TestNewSlidingWindow_ZeroConfigUsesDefaults(t *testing.T) {
	w := NewSlidingWindow(WindowConfig{})
	if w.cfg.Size != time.Minute {
		t.Errorf("expected default Size, got %v", w.cfg.Size)
	}
	if w.cfg.BucketCount != 10 {
		t.Errorf("expected default BucketCount, got %d", w.cfg.BucketCount)
	}
}

func TestSlidingWindow_InitialCountIsZero(t *testing.T) {
	w := NewSlidingWindow(DefaultWindowConfig())
	if c := w.Count(); c != 0 {
		t.Errorf("expected 0, got %d", c)
	}
}

func TestSlidingWindow_AddAndCount(t *testing.T) {
	w := NewSlidingWindow(DefaultWindowConfig())
	w.Add(3)
	w.Add(5)
	if c := w.Count(); c != 8 {
		t.Errorf("expected 8, got %d", c)
	}
}

func TestSlidingWindow_EvictsExpiredBuckets(t *testing.T) {
	base := time.Now()
	w := NewSlidingWindow(WindowConfig{Size: 10 * time.Second, BucketCount: 5})

	// inject a fake clock
	w.now = func() time.Time { return base }
	w.Add(10)

	// advance clock past the window
	w.now = func() time.Time { return base.Add(11 * time.Second) }
	w.Add(2)

	if c := w.Count(); c != 2 {
		t.Errorf("expected 2 after eviction, got %d", c)
	}
}

func TestSlidingWindow_Reset_ClearsAll(t *testing.T) {
	w := NewSlidingWindow(DefaultWindowConfig())
	w.Add(7)
	w.Reset()
	if c := w.Count(); c != 0 {
		t.Errorf("expected 0 after reset, got %d", c)
	}
}

func TestSlidingWindow_CountDoesNotEvictFutureEntries(t *testing.T) {
	base := time.Now()
	w := NewSlidingWindow(WindowConfig{Size: time.Minute, BucketCount: 5})
	w.now = func() time.Time { return base }
	w.Add(4)
	w.Add(6)

	// time has not advanced — nothing should be evicted
	if c := w.Count(); c != 10 {
		t.Errorf("expected 10, got %d", c)
	}
}
