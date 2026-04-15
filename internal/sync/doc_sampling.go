// Package sync provides synchronisation primitives for vaultpull.
//
// # Sampler
//
// Sampler decides whether a given profile sync operation should be
// "sampled" — i.e. selected for additional processing such as tracing,
// shadow writes, or detailed audit logging.
//
// Three modes are supported:
//
//   - SamplingModeAlways  – every profile is sampled.
//   - SamplingModeNever   – no profile is sampled.
//   - SamplingModeRandom  – each profile is sampled independently with
//     probability equal to Rate (0.0–1.0).
//
// Usage:
//
//	cfg := sync.DefaultSamplingConfig()
//	cfg.Rate = 0.25 // sample 25 % of syncs
//	sampler := sync.NewSampler(cfg)
//
//	if sampler.Sample(profile.Name) {
//	    // attach extra instrumentation
//	}
package sync
