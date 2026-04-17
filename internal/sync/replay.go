package sync

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// replayEntry is the on-disk representation of a saved secret payload.
type replayEntry struct {
	SavedAt time.Time         `json:"saved_at"`
	Secrets map[string]string `json:"secrets"`
}

// ReplayStore persists the last known secret payload per profile so it can be
// replayed when Vault is unreachable.
type ReplayStore struct {
	mu  sync.Mutex
	dir string
}

// NewReplayStore returns a ReplayStore that persists entries under dir.
func NewReplayStore(dir string) *ReplayStore {
	return &ReplayStore{dir: dir}
}

func (r *ReplayStore) path(profile string) string {
	return filepath.Join(r.dir, profile+".replay.json")
}

// Save persists secrets for the named profile.
func (r *ReplayStore) Save(profile string, secrets map[string]string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := os.MkdirAll(r.dir, 0o700); err != nil {
		return err
	}

	entry := replayEntry{SavedAt: time.Now().UTC(), Secrets: secrets}
	data, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	return os.WriteFile(r.path(profile), data, 0o600)
}

// Load retrieves the last saved payload for profile.
// Returns nil, false when no replay file exists.
func (r *ReplayStore) Load(profile string) (map[string]string, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	data, err := os.ReadFile(r.path(profile))
	if err != nil {
		return nil, false
	}

	var entry replayEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		return nil, false
	}
	return entry.Secrets, true
}

// Delete removes the replay file for profile.
func (r *ReplayStore) Delete(profile string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	err := os.Remove(r.path(profile))
	if os.IsNotExist(err) {
		return nil
	}
	return err
}
