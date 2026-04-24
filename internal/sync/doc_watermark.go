// Package sync provides synchronisation primitives for vaultpull.
//
// # High-Water Mark
//
// WatermarkTracker records the highest (most-recent) sync timestamp seen
// for each profile and exposes helpers to determine whether a candidate
// timestamp would advance the mark.
//
// Usage:
//
//	wm := sync.NewWatermarkTracker()
//	wm.Record("prod", time.Now())
//	if wm.IsHigher("prod", candidate) {
//	    // proceed with write
//	}
package sync
