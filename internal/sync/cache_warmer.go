package sync

import (
	"context"
	"fmt"
	"time"

	"github.com/yourusername/vaultpull/internal/config"
	"github.com/yourusername/vaultpull/internal/token"
)

// CacheWarmer pre-fetches and refreshes token cache entries for all profiles
// before a sync run, reducing latency during the actual sync.
type CacheWarmer struct {
	cfg    *config.Config
	cache  *token.Cache
	timeout time.Duration
}

// WarmResult holds the outcome of warming a single profile.
type WarmResult struct {
	Profile string
	Hit     bool
	Err     error
}

// NewCacheWarmer creates a CacheWarmer with the given config and token cache.
// timeout controls how long warming a single profile may take.
func NewCacheWarmer(cfg *config.Config, cache *token.Cache, timeout time.Duration) *CacheWarmer {
	if timeout <= 0 {
		timeout = 10 * time.Second
	}
	return &CacheWarmer{cfg: cfg, cache: cache, timeout: timeout}
}

// WarmAll iterates over every profile in the config and attempts to load its
// cached token. Results are returned for all profiles; a cache miss is not
// treated as an error — it simply means the token will be fetched at sync time.
func (w *CacheWarmer) WarmAll(ctx context.Context) []WarmResult {
	results := make([]WarmResult, 0, len(w.cfg.Profiles))
	for _, p := range w.cfg.Profiles {
		results = append(results, w.warmOne(ctx, p.Name))
	}
	return results
}

func (w *CacheWarmer) warmOne(ctx context.Context, profileName string) WarmResult {
	ctx, cancel := context.WithTimeout(ctx, w.timeout)
	defer cancel()

	doneCh := make(chan WarmResult, 1)
	go func() {
		entry, err := w.cache.Load(profileName)
		if err != nil {
			doneCh <- WarmResult{Profile: profileName, Hit: false, Err: err}
			return
		}
		doneCh <- WarmResult{Profile: profileName, Hit: entry != nil}
	}()

	select {
	case res := <-doneCh:
		return res
	case <-ctx.Done():
		return WarmResult{
			Profile: profileName,
			Hit:     false,
			Err:     fmt.Errorf("cache warm timed out for profile %q: %w", profileName, ctx.Err()),
		}
	}
}
