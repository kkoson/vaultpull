package sync

import "fmt"

// WithCachePolicy wraps a SyncFunc so that results are served from the
// in-memory CachePolicy when fresh, and written back when a fetch occurs in
// WriteThrough mode.
//
// The wrapped function signature matches the stage function type used by
// Pipeline: func(profile string) error.
func WithCachePolicy(policy *CachePolicy, fetch func(profile string) (map[string]string, error), consume func(profile string, secrets map[string]string) error) func(string) error {
	if policy == nil {
		panic("WithCachePolicy: policy must not be nil")
	}
	return func(profile string) error {
		if policy.IsFresh(profile) {
			secrets := policy.Get(profile)
			if secrets == nil {
				return fmt.Errorf("cache policy: fresh entry missing for profile %q", profile)
			}
			return consume(profile, secrets)
		}

		secrets, err := fetch(profile)
		if err != nil {
			return fmt.Errorf("cache policy: fetch failed for profile %q: %w", profile, err)
		}

		if policy.cfg.Mode == CachePolicyWriteThrough || policy.cfg.Mode == CachePolicyReadThrough {
			policy.Set(profile, secrets)
		}

		return consume(profile, secrets)
	}
}
