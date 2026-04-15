package sync

import (
	"testing"
	"time"
)

func TestDefaultShedderConfig_Values(t *testing.T) {
	cfg := DefaultShedderConfig()
	if cfg.MaxLoad != 100 {
		t.Errorf("expected MaxLoad 100, got %d", cfg.MaxLoad)
	}
	if cfg.WindowSize != 5*time.Second {
		t.Errorf("expected WindowSize 5s, got %v", cfg.WindowSize)
	}
	if cfg.DropPercent != 0.5 {
		t.Errorf("expected DropPercent 0.5, got %f", cfg.DropPercent)
	}
}

func TestNewLoadShedder_ZeroFallsBackToDefaults(t *testing.T) {
	s := NewLoadShedder(ShedderConfig{})
	if s.cfg.MaxLoad != 100 {
		t.Errorf("expected default MaxLoad 100, got %d", s.cfg.MaxLoad)
	}
}

func TestLoadShedder_AdmitsUnderLimit(t *testing.T) {
	s := NewLoadShedder(ShedderConfig{MaxLoad: 3, WindowSize: time.Second, DropPercent: 0.5})
	for i := 0; i < 3; i++ {
		if !s.Admit() {
			t.Fatalf("expected admission on call %d", i+1)
		}
	}
}

func TestLoadShedder_ShedsAtLimit(t *testing.T) {
	s := NewLoadShedder(ShedderConfig{MaxLoad: 2, WindowSize: time.Second, DropPercent: 0.5})
	s.Admit()
	s.Admit()
	if s.Admit() {
		t.Error("expected shed when at MaxLoad")
	}
}

func TestLoadShedder_Release_DecrementsCounter(t *testing.T) {
	s := NewLoadShedder(ShedderConfig{MaxLoad: 1, WindowSize: time.Second, DropPercent: 0.5})
	s.Admit()
	s.Release()
	if !s.Admit() {
		t.Error("expected admission after release")
	}
}

func TestLoadShedder_Release_DoesNotGoBelowZero(t *testing.T) {
	s := NewLoadShedder(ShedderConfig{MaxLoad: 5, WindowSize: time.Second, DropPercent: 0.5})
	s.Release() // should not panic or underflow
	current, _, _ := s.Stats()
	if current != 0 {
		t.Errorf("expected current 0, got %d", current)
	}
}

func TestLoadShedder_Stats_Counters(t *testing.T) {
	s := NewLoadShedder(ShedderConfig{MaxLoad: 1, WindowSize: time.Second, DropPercent: 0.5})
	s.Admit()
	s.Admit() // shed
	current, shed, total := s.Stats()
	if current != 1 {
		t.Errorf("expected current 1, got %d", current)
	}
	if shed != 1 {
		t.Errorf("expected shed 1, got %d", shed)
	}
	if total != 2 {
		t.Errorf("expected total 2, got %d", total)
	}
}

func TestLoadShedder_WindowReset(t *testing.T) {
	s := NewLoadShedder(ShedderConfig{MaxLoad: 1, WindowSize: 10 * time.Millisecond, DropPercent: 0.5})
	s.Admit()
	time.Sleep(20 * time.Millisecond)
	if !s.Admit() {
		t.Error("expected admission after window reset")
	}
	_, _, total := s.Stats()
	if total != 1 {
		t.Errorf("expected total reset to 1, got %d", total)
	}
}
