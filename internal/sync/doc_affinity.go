// Package sync provides synchronisation primitives for vaultpull.
//
// # Affinity
//
// AffinityRouter assigns profiles to named execution slots ("affinities").
// Profiles tagged with the same affinity key are guaranteed to run on the
// same logical worker, which is useful when downstream systems require
// ordered or serialised writes for a given environment.
//
// Usage:
//
//	router := sync.NewAffinityRouter(sync.DefaultAffinityConfig)
//	slot := router.Assign(profile)
//	workerID := slot % numWorkers
//
// If no affinity tag is present on the profile the router falls back to
// a hash of the profile name so load is still spread evenly.
package sync
