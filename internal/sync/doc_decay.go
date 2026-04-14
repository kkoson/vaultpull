// Package sync provides synchronisation primitives and helpers used by
// the vaultpull sync engine.
//
// # Decay Counter
//
// DecayCounter implements an exponentially decaying counter that models
// a time-weighted signal. Each call to Add first ages the current value
// by the fraction of a half-life that has elapsed since the last write,
// then adds the new delta.
//
// This is useful for tracking error rates or request frequencies where
// recent events should carry more weight than older ones.
//
//	dc := sync.NewDecayCounter(sync.DefaultDecayConfig())
//	dc.Add(1)              // record an event
//	v := dc.Value()        // read the current decayed value
package sync
