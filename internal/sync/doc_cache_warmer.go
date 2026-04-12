// Package sync provides synchronisation primitives and orchestration for
// pulling secrets from HashiCorp Vault into local .env files.
//
// # Cache Warmer
//
// CacheWarmer pre-fetches token-cache entries for every configured profile
// before a sync run begins.  Warming the cache reduces per-profile latency
// because the token is already resident in memory by the time the Vault client
// needs it.
//
// Usage:
//
//	 warmer := sync.NewCacheWarmer(cfg, cache, 5*time.Second)
//	 results := warmer.WarmAll(ctx)
//	 for _, r := range results {
//	     if r.Err != nil {
//	         log.Printf("warm failed for %s: %v", r.Profile, r.Err)
//	     }
//	 }
//
// A cache miss (Hit == false, Err == nil) is not an error condition; it simply
// means the token will be fetched fresh during the sync.
package sync
