package sync

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// CheckpointEntry records the last successful sync for a profile.
type CheckpointEntry struct {
	Profile   string    `json:"profile"`
	SyncedAt  time.Time `json:"synced_at"`
	SecretPath string   `json:"secret_path"`
}

// Checkpoint persists sync state between runs so that profiles
// can be skipped when nothing has changed.
type Checkpoint struct {
	mu      sync.RWMutex
	path    string
	entries map[string]CheckpointEntry
}

// NewCheckpoint loads (or creates) a checkpoint file at the given path.
func NewCheckpoint(path string) (*Checkpoint, error) {
	cp := &Checkpoint{
		path:    path,
		entries: make(map[string]CheckpointEntry),
	}
	if err := cp.load(); err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	return cp, nil
}

// Get returns the checkpoint entry for a profile, and whether it exists.
func (c *Checkpoint) Get(profile string) (CheckpointEntry, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	e, ok := c.entries[profile]
	return e, ok
}

// Set records a successful sync for a profile and flushes to disk.
func (c *Checkpoint) Set(entry CheckpointEntry) error {
	c.mu.Lock()
	c.entries[entry.Profile] = entry
	c.mu.Unlock()
	return c.flush()
}

// Clear removes the checkpoint entry for a profile.
func (c *Checkpoint) Clear(profile string) error {
	c.mu.Lock()
	delete(c.entries, profile)
	c.mu.Unlock()
	return c.flush()
}

func (c *Checkpoint) load() error {
	data, err := os.ReadFile(c.path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &c.entries)
}

func (c *Checkpoint) flush() error {
	if err := os.MkdirAll(filepath.Dir(c.path), 0o700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(c.entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(c.path, data, 0o600)
}
