package sync

import (
	"errors"
	"sync"
	"time"
)

// FenceConfig holds configuration for the write fence.
type FenceConfig struct {
	// Window is the duration during which repeated writes to the same key are blocked.
	Window time.Duration
}

// DefaultFenceConfig returns sensible defaults for the write fence.
func DefaultFenceConfig() FenceConfig {
	return FenceConfig{
		Window: 30 * time.Second,
	}
}

// WriteFence prevents duplicate writes to the same env key within a time window.
type WriteFence struct {
	cfg    FenceConfig
	mu     sync.Mutex
	seen   map[string]time.Time
	nowFn  func() time.Time
}

// NewWriteFence creates a new WriteFence with the given config.
// Zero-value fields fall back to defaults.
func NewWriteFence(cfg FenceConfig) *WriteFence {
	if cfg.Window <= 0 {
		cfg.Window = DefaultFenceConfig().Window
	}
	return &WriteFence{
		cfg:   cfg,
		seen:  make(map[string]time.Time),
		nowFn: time.Now,
	}
}

// Allow returns nil if the key may be written, or an error if it is fenced.
func (f *WriteFence) Allow(key string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	now := f.nowFn()
	if t, ok := f.seen[key]; ok && now.Sub(t) < f.cfg.Window {
		return errors.New("write fenced: key " + key + " was written recently")
	}
	f.seen[key] = now
	return nil
}

// Reset clears the fence record for a specific key.
func (f *WriteFence) Reset(key string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	delete(f.seen, key)
}

// Flush clears all fence records.
func (f *WriteFence) Flush() {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.seen = make(map[string]time.Time)
}
