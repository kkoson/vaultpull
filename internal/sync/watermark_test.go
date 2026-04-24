package sync

import (
	"testing"
	"time"
)

func TestNewWatermarkTracker_NotNil(t *testing.T) {
	wm := NewWatermarkTracker()
	if wm == nil {
		t.Fatal("expected non-nil WatermarkTracker")
	}
}

func TestWatermarkTracker_Len_Empty(t *testing.T) {
	wm := NewWatermarkTracker()
	if wm.Len() != 0 {
		t.Fatalf("expected 0, got %d", wm.Len())
	}
}

func TestWatermarkTracker_Get_Miss(t *testing.T) {
	wm := NewWatermarkTracker()
	_, ok := wm.Get("missing")
	if ok {
		t.Fatal("expected miss for unknown profile")
	}
}

func TestWatermarkTracker_RecordAndGet(t *testing.T) {
	wm := NewWatermarkTracker()
	now := time.Now()
	wm.Record("prod", now)
	got, ok := wm.Get("prod")
	if !ok {
		t.Fatal("expected hit after Record")
	}
	if !got.Equal(now) {
		t.Fatalf("expected %v, got %v", now, got)
	}
}

func TestWatermarkTracker_Record_OnlyAdvances(t *testing.T) {
	wm := NewWatermarkTracker()
	now := time.Now()
	older := now.Add(-time.Minute)
	wm.Record("prod", now)
	wm.Record("prod", older) // should be ignored
	got, _ := wm.Get("prod")
	if !got.Equal(now) {
		t.Fatalf("mark should not regress: got %v", got)
	}
}

func TestWatermarkTracker_IsHigher_NoMark(t *testing.T) {
	wm := NewWatermarkTracker()
	if !wm.IsHigher("prod", time.Now()) {
		t.Fatal("IsHigher should return true when no mark recorded")
	}
}

func TestWatermarkTracker_IsHigher_FutureTimestamp(t *testing.T) {
	wm := NewWatermarkTracker()
	now := time.Now()
	wm.Record("prod", now)
	if !wm.IsHigher("prod", now.Add(time.Second)) {
		t.Fatal("expected true for timestamp after mark")
	}
}

func TestWatermarkTracker_IsHigher_PastTimestamp(t *testing.T) {
	wm := NewWatermarkTracker()
	now := time.Now()
	wm.Record("prod", now)
	if wm.IsHigher("prod", now.Add(-time.Second)) {
		t.Fatal("expected false for timestamp before mark")
	}
}

func TestWatermarkTracker_Reset(t *testing.T) {
	wm := NewWatermarkTracker()
	wm.Record("prod", time.Now())
	wm.Reset("prod")
	_, ok := wm.Get("prod")
	if ok {
		t.Fatal("expected miss after Reset")
	}
	if wm.Len() != 0 {
		t.Fatalf("expected Len 0 after Reset, got %d", wm.Len())
	}
}
