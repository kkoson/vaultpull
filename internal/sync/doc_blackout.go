// Package sync provides synchronisation primitives for vaultpull.
//
// # Blackout Windows
//
// BlackoutManager prevents sync operations from running during
// configured time windows (e.g. maintenance periods or business-
// critical hours).
//
// Usage:
//
//	cfg := sync.DefaultBlackoutConfig()
//	cfg.Windows = []sync.BlackoutWindow{
//		{Start: "22:00", End: "06:00"},
//	}
//	bm := sync.NewBlackoutManager(cfg)
//
// Use WithBlackout to wrap a pipeline stage so it skips execution
// automatically during any configured window.
package sync
