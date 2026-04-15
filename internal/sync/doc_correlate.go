// Package sync provides synchronisation primitives for vaultpull.
//
// # Correlator
//
// Correlator attaches a unique correlation ID to every sync operation so that
// log lines, audit entries, and trace spans produced by a single profile run
// can be grouped together during post-hoc analysis.
//
// Usage:
//
//	corr := sync.NewCorrelator()
//	ctx  := corr.Inject(ctx, "my-profile")
//	id   := sync.CorrelationIDFromContext(ctx) // stable for the lifetime of the run
package sync
