package sync

import (
	"sync"
	"time"
)

// WindowConfig holds configuration for a sliding window counter.
type WindowConfig struct {
	// Size is the duration of the sliding window.
	Size time.Duration
	// BucketCount is the number of buckets used to approximate the window.
	BucketCount int
}

// DefaultWindowConfig returns a WindowConfig with sensible defaults.
func DefaultWindowConfig() WindowConfig {
	return WindowConfig{
		Size:        time.Minute,
		BucketCount: 10,
	}
}

// bucket holds a count and the time it was created.
type bucket struct {
	count int
	at    time.Time
}

// SlidingWindow is a thread-safe sliding window counter.
type SlidingWindow struct {
	mu      sync.Mutex
	cfg     WindowConfig
	buckets []bucket
	now     func() time.Time
}

// NewSlidingWindow creates a new SlidingWindow with the given config.
// If cfg.BucketCount <= 0 or cfg.Size <= 0, defaults are applied.
func NewSlidingWindow(cfg WindowConfig) *SlidingWindow {
	def := DefaultWindowConfig()
	if cfg.Size <= 0 {
		cfg.Size = def.Size
	}
	if cfg.BucketCount <= 0 {
		cfg.BucketCount = def.BucketCount
	}
	return &SlidingWindow{
		cfg:     cfg,
		buckets: make([]bucket, 0, cfg.BucketCount),
		now:     time.Now,
	}
}

// Add increments the window counter by n.
func (w *SlidingWindow) Add(n int) {
	w.mu.Lock()
	defer w.mu.Unlock()
	now := w.now()
	w.evict(now)
	w.buckets = append(w.buckets, bucket{count: n, at: now})
}

// Count returns the total count within the current window.
func (w *SlidingWindow) Count() int {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.evict(w.now())
	total := 0
	for _, b := range w.buckets {
		total += b.count
	}
	return total
}

// Reset clears all buckets.
func (w *SlidingWindow) Reset() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.buckets = w.buckets[:0]
}

// evict removes buckets older than the window size. Must be called with mu held.
func (w *SlidingWindow) evict(now time.Time) {
	cutoff := now.Add(-w.cfg.Size)
	i := 0
	for i < len(w.buckets) && w.buckets[i].at.Before(cutoff) {
		i++
	}
	w.buckets = w.buckets[i:]
}
