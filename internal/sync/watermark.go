package sync

import (
	"sync"
	"time"
)

// WatermarkTracker tracks the highest observed sync timestamp per profile.
type WatermarkTracker struct {
	mu    sync.RWMutex
	marks map[string]time.Time
}

// NewWatermarkTracker returns an initialised WatermarkTracker.
func NewWatermarkTracker() *WatermarkTracker {
	return &WatermarkTracker{
		marks: make(map[string]time.Time),
	}
}

// Record stores t as the high-water mark for profile if t is later than the
// current mark (or if no mark has been recorded yet).
func (w *WatermarkTracker) Record(profile string, t time.Time) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if cur, ok := w.marks[profile]; !ok || t.After(cur) {
		w.marks[profile] = t
	}
}

// Get returns the current high-water mark for profile and whether one exists.
func (w *WatermarkTracker) Get(profile string) (time.Time, bool) {
	w.mu.RLock()
	defer w.mu.RUnlock()
	t, ok := w.marks[profile]
	return t, ok
}

// IsHigher reports whether t is strictly after the current high-water mark for
// profile. Returns true when no mark has been recorded yet.
func (w *WatermarkTracker) IsHigher(profile string, t time.Time) bool {
	w.mu.RLock()
	defer w.mu.RUnlock()
	cur, ok := w.marks[profile]
	if !ok {
		return true
	}
	return t.After(cur)
}

// Reset removes the high-water mark for profile.
func (w *WatermarkTracker) Reset(profile string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	delete(w.marks, profile)
}

// Len returns the number of profiles currently tracked.
func (w *WatermarkTracker) Len() int {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return len(w.marks)
}
