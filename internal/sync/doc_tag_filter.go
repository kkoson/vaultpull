// Package sync provides profile-level tag filtering via TagFilter.
//
// TagFilter allows callers to restrict which profiles participate in a sync
// run based on string labels (tags) attached to each profile in the
// configuration file.
//
// Usage:
//
//	// Only sync profiles tagged "prod" or "critical".
//	f := sync.NewTagFilter([]string{"prod", "critical"})
//
//	for _, p := range cfg.Profiles {
//		if !f.Allow(p.Tags) {
//			continue
//		}
//		// … run syncer for p
//	}
//
// An empty tag list makes the filter a no-op: every profile is allowed.
package sync
