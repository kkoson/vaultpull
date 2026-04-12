// Package sync provides synchronisation primitives for vaultpull.
//
// # Tracer
//
// Tracer records fine-grained timing events that occur during a sync run.
// It supports two verbosity levels:
//
//   - TraceLevelBasic – emits one event per profile (start / end).
//   - TraceLevelFull  – additionally emits per-stage timings produced by
//     pipeline middleware such as WithTiming.
//
// Usage:
//
//	tr := sync.NewTracer(sync.TraceLevelFull, os.Stderr)
//	tr.Record("prod", "validate", "done", 2*time.Millisecond)
//	entries := tr.Entries()
package sync
