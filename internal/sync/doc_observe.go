// Package sync provides the Observer type for recording and optionally
// printing per-profile and per-key sync events.
//
// # Observer
//
// An Observer collects ObserveEvent records produced during a sync run.
// Three verbosity levels are supported:
//
//   - ObserveOff    – no recording or output
//   - ObserveSummary – one line per profile (key is empty)
//   - ObserveFull    – one line per key-level change
//
// # Middleware
//
// WithObserver wraps a pipeline StageFunc, recording a "synced" or "failed"
// summary event after each profile completes. It also stores the observer
// in the context so downstream stages can call ObserverFromContext to emit
// fine-grained key events.
//
// Example:
//
//	obs := sync.NewObserver(sync.ObserveFull, os.Stdout)
//	pipeline.AddStage(sync.WithObserver(obs)(myStage))
package sync
