// Package sync – cache policy
//
// CachePolicy provides an in-process, TTL-based caching layer that sits in
// front of Vault secret fetches. Three modes are supported:
//
//   - CachePolicyBypass: every sync goes directly to Vault; the cache is
//     never consulted or written.
//
//   - CachePolicyReadThrough (default): if a fresh entry exists for the
//     profile the cached secrets are returned immediately; otherwise Vault
//     is queried and the result is stored for subsequent requests.
//
//   - CachePolicyWriteThrough: Vault is always queried but the result is
//     written back to the in-memory store so that future reads within the
//     TTL window are served without a network round-trip.
//
// Usage:
//
//	policy := sync.NewCachePolicy(sync.DefaultCachePolicyConfig())
//	stage := sync.WithCachePolicy(policy, fetchFn, consumeFn)
//	// stage can be added to a Pipeline or called directly.
package sync
