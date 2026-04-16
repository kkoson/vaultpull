package sync

import (
	"testing"
	"time"
)

func fixedNow(h, m int) func() time.Time {
	return func() time.Time {
		base := time.Date(2024, 1, 1, h, m, 0, 0, time.UTC)
		return base
	}
}

func TestDefaultBlackoutConfig_NoWindows(t *testing.T) {
	cfg := DefaultBlackoutConfig()
	if len(cfg.Windows) != 0 {
		t.Fatalf("expected no windows, got %d", len(cfg.Windows))
	}
}

func TestBlackoutManager_NotBlackedOut_NoWindows(t *testing.T) {
	bm := NewBlackoutManager(DefaultBlackoutConfig())
	if bm.IsBlackedOut() {
		t.Fatal("expected not blacked out with no windows")
	}
}

func TestBlackoutManager_BlackedOut_WithinWindow(t *testing.T) {
	cfg := BlackoutConfig{
		Windows: []BlackoutWindow{
			{Start: 2 * time.Hour, End: 4 * time.Hour},
		},
	}
	bm := NewBlackoutManager(cfg)
	bm.nowFunc = fixedNow(3, 0)
	if !bm.IsBlackedOut() {
		t.Fatal("expected blacked out at 03:00")
	}
}

func TestBlackoutManager_NotBlackedOut_OutsideWindow(t *testing.T) {
	cfg := BlackoutConfig{
		Windows: []BlackoutWindow{
			{Start: 2 * time.Hour, End: 4 * time.Hour},
		},
	}
	bm := NewBlackoutManager(cfg)
	bm.nowFunc = fixedNow(5, 0)
	if bm.IsBlackedOut() {
		t.Fatal("expected not blacked out at 05:00")
	}
}

func TestBlackoutManager_WrapsAroundMidnight(t *testing.T) {
	cfg := BlackoutConfig{
		Windows: []BlackoutWindow{
			{Start: 23 * time.Hour, End: 1 * time.Hour},
		},
	}
	bm := NewBlackoutManager(cfg)
	bm.nowFunc = fixedNow(23, 30)
	if !bm.IsBlackedOut() {
		t.Fatal("expected blacked out at 23:30 in wrap-around window")
	}
	bm.nowFunc = fixedNow(0, 30)
	if !bm.IsBlackedOut() {
		t.Fatal("expected blacked out at 00:30 in wrap-around window")
	}
}

func TestBlackoutManager_AddWindow(t *testing.T) {
	bm := NewBlackoutManager(DefaultBlackoutConfig())
	bm.nowFunc = fixedNow(10, 0)
	if bm.IsBlackedOut() {
		t.Fatal("expected not blacked out before adding window")
	}
	bm.AddWindow(BlackoutWindow{Start: 9 * time.Hour, End: 11 * time.Hour})
	if !bm.IsBlackedOut() {
		t.Fatal("expected blacked out after adding window")
	}
}
