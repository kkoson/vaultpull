package sync

import (
	"strings"
	"testing"
)

func TestDefaultBudgetConfig_Values(t *testing.T) {
	cfg := DefaultBudgetConfig()
	if cfg.MaxErrorRate != 0.5 {
		t.Errorf("expected MaxErrorRate 0.5, got %v", cfg.MaxErrorRate)
	}
	if cfg.WindowSize != 10 {
		t.Errorf("expected WindowSize 10, got %d", cfg.WindowSize)
	}
	if cfg.MinSampleSize != 3 {
		t.Errorf("expected MinSampleSize 3, got %d", cfg.MinSampleSize)
	}
}

func TestNewErrorBudget_ZeroConfigUsesDefaults(t *testing.T) {
	b := NewErrorBudget(BudgetConfig{})
	def := DefaultBudgetConfig()
	if b.cfg.WindowSize != def.WindowSize {
		t.Errorf("expected window %d, got %d", def.WindowSize, b.cfg.WindowSize)
	}
}

func TestErrorBudget_NotExhausted_BelowMinSample(t *testing.T) {
	b := NewErrorBudget(BudgetConfig{MaxErrorRate: 0.5, WindowSize: 10, MinSampleSize: 3})
	b.Record(false)
	b.Record(false)
	if b.Exhausted() {
		t.Error("budget should not be exhausted below MinSampleSize")
	}
}

func TestErrorBudget_NotExhausted_LowErrorRate(t *testing.T) {
	b := NewErrorBudget(BudgetConfig{MaxErrorRate: 0.5, WindowSize: 10, MinSampleSize: 3})
	for i := 0; i < 8; i++ {
		b.Record(true)
	}
	b.Record(false)
	b.Record(false)
	if b.Exhausted() {
		t.Error("budget should not be exhausted at 20%% error rate")
	}
}

func TestErrorBudget_Exhausted_HighErrorRate(t *testing.T) {
	b := NewErrorBudget(BudgetConfig{MaxErrorRate: 0.5, WindowSize: 10, MinSampleSize: 3})
	for i := 0; i < 4; i++ {
		b.Record(false)
	}
	for i := 0; i < 2; i++ {
		b.Record(true)
	}
	if !b.Exhausted() {
		t.Error("budget should be exhausted at >50%% error rate")
	}
}

func TestErrorBudget_SlidingWindow_OldFailuresDropped(t *testing.T) {
	b := NewErrorBudget(BudgetConfig{MaxErrorRate: 0.5, WindowSize: 4, MinSampleSize: 4})
	// Fill with failures.
	for i := 0; i < 4; i++ {
		b.Record(false)
	}
	// Overwrite all slots with successes.
	for i := 0; i < 4; i++ {
		b.Record(true)
	}
	if b.Exhausted() {
		t.Error("old failures should have been overwritten by successes")
	}
}

func TestErrorBudget_Stats_NoSamples(t *testing.T) {
	b := NewErrorBudget(DefaultBudgetConfig())
	if !strings.Contains(b.Stats(), "no samples") {
		t.Errorf("expected 'no samples' in stats, got: %s", b.Stats())
	}
}

func TestErrorBudget_Stats_WithSamples(t *testing.T) {
	b := NewErrorBudget(BudgetConfig{MaxErrorRate: 0.5, WindowSize: 4, MinSampleSize: 2})
	b.Record(true)
	b.Record(false)
	s := b.Stats()
	if !strings.Contains(s, "budget:") {
		t.Errorf("unexpected stats format: %s", s)
	}
}
