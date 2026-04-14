package sync

import (
	"fmt"
	"sync"
	"time"
)

// DefaultBudgetConfig returns a BudgetConfig with sensible defaults.
func DefaultBudgetConfig() BudgetConfig {
	return BudgetConfig{
		MaxErrorRate:  0.5,
		WindowSize:    10,
		MinSampleSize: 3,
	}
}

// BudgetConfig controls error-budget enforcement behaviour.
type BudgetConfig struct {
	// MaxErrorRate is the fraction of failures allowed before the budget is
	// exhausted (0.0–1.0).
	MaxErrorRate float64
	// WindowSize is the number of recent outcomes tracked.
	WindowSize int
	// MinSampleSize is the minimum number of outcomes required before the
	// budget can be considered exhausted.
	MinSampleSize int
}

// ErrorBudget tracks a sliding window of success/failure outcomes and reports
// whether the error budget has been exhausted.
type ErrorBudget struct {
	mu      sync.Mutex
	cfg     BudgetConfig
	window  []bool // true == success
	pos     int
	full    bool
	updated time.Time
}

// NewErrorBudget creates an ErrorBudget using cfg. Zero-value fields fall back
// to DefaultBudgetConfig.
func NewErrorBudget(cfg BudgetConfig) *ErrorBudget {
	def := DefaultBudgetConfig()
	if cfg.MaxErrorRate <= 0 {
		cfg.MaxErrorRate = def.MaxErrorRate
	}
	if cfg.WindowSize <= 0 {
		cfg.WindowSize = def.WindowSize
	}
	if cfg.MinSampleSize <= 0 {
		cfg.MinSampleSize = def.MinSampleSize
	}
	return &ErrorBudget{
		cfg:    cfg,
		window: make([]bool, cfg.WindowSize),
	}
}

// Record registers a single outcome. success=true means the operation
// succeeded.
func (b *ErrorBudget) Record(success bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.window[b.pos] = success
	b.pos = (b.pos + 1) % b.cfg.WindowSize
	if b.pos == 0 {
		b.full = true
	}
	b.updated = time.Now()
}

// Exhausted returns true when the error rate across the current window exceeds
// MaxErrorRate and at least MinSampleSize outcomes have been recorded.
func (b *ErrorBudget) Exhausted() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	size := b.cfg.WindowSize
	if !b.full {
		size = b.pos
	}
	if size < b.cfg.MinSampleSize {
		return false
	}
	failures := 0
	for i := 0; i < size; i++ {
		if !b.window[i] {
			failures++
		}
	}
	return float64(failures)/float64(size) > b.cfg.MaxErrorRate
}

// Stats returns a human-readable summary of the current window.
func (b *ErrorBudget) Stats() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	size := b.cfg.WindowSize
	if !b.full {
		size = b.pos
	}
	if size == 0 {
		return "budget: no samples"
	}
	failures := 0
	for i := 0; i < size; i++ {
		if !b.window[i] {
			failures++
		}
	}
	return fmt.Sprintf("budget: %d/%d failures (%.0f%%)",
		failures, size, float64(failures)/float64(size)*100)
}
