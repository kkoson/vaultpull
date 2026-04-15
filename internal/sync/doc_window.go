// Package sync provides synchronisation primitives and helpers for vaultpull.
//
// # Sliding Window
//
// SlidingWindow is a thread-safe approximate sliding-window counter.
// It divides time into discrete buckets and evicts buckets that fall
// outside the configured window on every read or write, giving an
// O(n) space-bounded approximation of a true sliding window.
//
// Typical usage:
//
//	w := sync.NewSlidingWindow(sync.DefaultWindowConfig())
//	w.Add(1)              // record an event
//	count := w.Count()   // total events in the last window
//	w.Reset()            // clear all buckets
package sync
