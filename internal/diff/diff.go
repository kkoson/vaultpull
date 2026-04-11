// Package diff provides utilities for comparing secret maps
// and reporting changes between vault secrets and local env files.
package diff

// ChangeType represents the type of change for a secret key.
type ChangeType int

const (
	Added ChangeType = iota
	Updated
	Removed
	Unchanged
)

// Change represents a single secret key change.
type Change struct {
	Key      string
	Type     ChangeType
	OldValue string
	NewValue string
}

// Result holds the full diff between two secret maps.
type Result struct {
	Changes []Change
}

// HasChanges returns true if there are any added, updated, or removed entries.
func (r *Result) HasChanges() bool {
	for _, c := range r.Changes {
		if c.Type != Unchanged {
			return true
		}
	}
	return false
}

// Summary returns counts of each change type.
func (r *Result) Summary() (added, updated, removed, unchanged int) {
	for _, c := range r.Changes {
		switch c.Type {
		case Added:
			added++
		case Updated:
			updated++
		case Removed:
			removed++
		case Unchanged:
			unchanged++
		}
	}
	return
}

// Compare computes the diff between existing local secrets and incoming vault secrets.
// existing is the current local env map; incoming is the fetched vault secret map.
func Compare(existing, incoming map[string]string) *Result {
	result := &Result{}

	for k, newVal := range incoming {
		if oldVal, ok := existing[k]; ok {
			if oldVal == newVal {
				result.Changes = append(result.Changes, Change{Key: k, Type: Unchanged, OldValue: oldVal, NewValue: newVal})
			} else {
				result.Changes = append(result.Changes, Change{Key: k, Type: Updated, OldValue: oldVal, NewValue: newVal})
			}
		} else {
			result.Changes = append(result.Changes, Change{Key: k, Type: Added, NewValue: newVal})
		}
	}

	for k, oldVal := range existing {
		if _, ok := incoming[k]; !ok {
			result.Changes = append(result.Changes, Change{Key: k, Type: Removed, OldValue: oldVal})
		}
	}

	return result
}
