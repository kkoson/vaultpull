package sync

import (
	"fmt"

	"github.com/owner/vaultpull/internal/diff"
)

// DriftResult holds the outcome of a drift check for a single profile.
type DriftResult struct {
	// Profile is the name of the profile that was checked.
	Profile string

	// Changes contains the diff entries between the snapshot and current secrets.
	Changes []diff.Entry

	// HasDrift is true when at least one change was detected.
	HasDrift bool
}

// DriftDetector compares a SnapshotStore against a current set of secrets
// and reports any differences.
type DriftDetector struct {
	store *SnapshotStore
}

// NewDriftDetector creates a DriftDetector backed by the given SnapshotStore.
func NewDriftDetector(store *SnapshotStore) (*DriftDetector, error) {
	if store == nil {
		return nil, fmt.Errorf("drift: snapshot store must not be nil")
	}
	return &DriftDetector{store: store}, nil
}

// Check compares the current secrets for a profile against the last saved
// snapshot. If no snapshot exists the result will report all keys as added.
func (d *DriftDetector) Check(profile string, current map[string]string) DriftResult {
	prev, _ := d.store.Get(profile)

	entries := diff.Compare(prev, current)

	hasChanges := false
	for _, e := range entries {
		if e.Kind != diff.KindUnchanged {
			hasChanges = true
			break
		}
	}

	return DriftResult{
		Profile:  profile,
		Changes:  entries,
		HasDrift: hasChanges,
	}
}
