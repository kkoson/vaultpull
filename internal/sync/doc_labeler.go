// Package sync — labeler
//
// Labeler assigns key/value label pairs to a sync profile based on
// configurable rules. Each LabelRule contains:
//
//   - Key:       the label key written to the output map.
//   - Value:     the label value written when the predicate matches.
//   - Predicate: a function(profile, tags) bool that controls matching.
//
// Two built-in predicates are provided:
//
//	HasTag(tag)        – matches when the profile carries the given tag.
//	ProfilePrefix(pfx) – matches when the profile name starts with pfx.
//
// Example:
//
//	l := sync.NewLabeler([]sync.LabelRule{
//		{Key: "env", Predicate: sync.ProfilePrefix("prod-"), Value: "production"},
//		{Key: "critical", Predicate: sync.HasTag("critical"), Value: "true"},
//	})
//	labels := l.Label("prod-api", []string{"critical"})
//	// labels == map[string]string{"env": "production", "critical": "true"}
package sync
