package sync

import "time"

// DriftDetector compares the current state of secrets against a previously
// saved snapshot to identify keys that have been added, removed, or changed
// outside of a normal vaultpull sync cycle.
type DriftDetector struct {
	store       *SnapshotStore
	lastChecked time.Time
}

// NewDriftDetector returns a DriftDetector backed by the given SnapshotStore.
func NewDriftDetector(store *SnapshotStore) *DriftDetector {
	return &DriftDetector{store: store}
}

// Detect compares current against the stored snapshot for profile.
// It returns (true, driftedKeys, nil) when drift is detected, or
// (false, nil, nil) when the state matches or no snapshot exists yet.
func (d *DriftDetector) Detect(profile string, current map[string]string) (bool, []string, error) {
	d.lastChecked = time.Now()

	snapshot, err := d.store.Get(profile)
	if err != nil {
		return false, nil, err
	}

	// No baseline yet — treat as clean so we don't false-positive on first run.
	if snapshot == nil {
		return false, nil, nil
	}

	var drifted []string

	// Keys that changed value or were removed.
	for k, snapVal := range snapshot {
		curVal, ok := current[k]
		if !ok {
			drifted = append(drifted, k)
			continue
		}
		if curVal != snapVal {
			drifted = append(drifted, k)
		}
	}

	// Keys that were added since the snapshot.
	for k := range current {
		if _, ok := snapshot[k]; !ok {
			drifted = append(drifted, k)
		}
	}

	if len(drifted) == 0 {
		return false, nil, nil
	}

	return true, drifted, nil
}

// LastChecked returns the timestamp of the most recent Detect call.
func (d *DriftDetector) LastChecked() time.Time {
	return d.lastChecked
}
