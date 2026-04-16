package sync

import (
	"sync"
	"time"
)

// BlackoutConfig holds configuration for the blackout window manager.
type BlackoutConfig struct {
	// Windows is a list of time ranges during which syncs are suppressed.
	Windows []BlackoutWindow
}

// BlackoutWindow defines a daily recurring time range.
type BlackoutWindow struct {
	Start time.Duration // offset from midnight
	End   time.Duration // offset from midnight
}

// DefaultBlackoutConfig returns a BlackoutConfig with no windows.
func DefaultBlackoutConfig() BlackoutConfig {
	return BlackoutConfig{}
}

// BlackoutManager suppresses operations during configured time windows.
type BlackoutManager struct {
	mu      sync.RWMutex
	config  BlackoutConfig
	nowFunc func() time.Time
}

// NewBlackoutManager creates a new BlackoutManager with the given config.
func NewBlackoutManager(cfg BlackoutConfig) *BlackoutManager {
	return &BlackoutManager{
		config:  cfg,
		nowFunc: time.Now,
	}
}

// IsBlackedOut returns true if the current time falls within any blackout window.
func (b *BlackoutManager) IsBlackedOut() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()

	now := b.nowFunc()
	midnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	offset := now.Sub(midnight)

	for _, w := range b.config.Windows {
		if w.Start <= w.End {
			if offset >= w.Start && offset < w.End {
				return true
			}
		} else {
			// wraps midnight
			if offset >= w.Start || offset < w.End {
				return true
			}
		}
	}
	return false
}

// AddWindow appends a new blackout window at runtime.
func (b *BlackoutManager) AddWindow(w BlackoutWindow) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.config.Windows = append(b.config.Windows, w)
}
