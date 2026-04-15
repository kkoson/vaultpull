// Package sync provides the WriteFence primitive, which prevents duplicate
// writes to the same output file within a configurable time window.
//
// # WriteFence
//
// A WriteFence tracks the last write timestamp for each output file path.
// If the same path is written again before the window elapses, the fence
// returns an error and the write is skipped.
//
// Usage:
//
//	fence := sync.NewWriteFence(sync.DefaultFenceConfig())
//	stage := sync.WithWriteFence(fence, mySyncFn)
//
// The middleware logs a message and returns nil (skips) rather than
// propagating the fence error, so callers treat fenced writes as no-ops.
package sync
