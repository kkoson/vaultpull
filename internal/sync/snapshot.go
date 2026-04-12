package sync

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

// Snapshot captures the state of secrets for a profile at a point in time.
type Snapshot struct {
	Profile   string            `json:"profile"`
	CapturedAt time.Time        `json:"captured_at"`
	Secrets   map[string]string `json:"secrets"`
}

// SnapshotStore persists and retrieves snapshots from disk.
type SnapshotStore struct {
	mu   sync.RWMutex
	path string
	data map[string]Snapshot
}

// NewSnapshotStore creates a SnapshotStore backed by the given file path.
func NewSnapshotStore(path string) (*SnapshotStore, error) {
	s := &SnapshotStore{
		path: path,
		data: make(map[string]Snapshot),
	}
	if err := s.load(); err != nil {
		return nil, err
	}
	return s, nil
}

// Save stores a snapshot for the given profile and flushes to disk.
func (s *SnapshotStore) Save(snap Snapshot) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[snap.Profile] = snap
	return s.flush()
}

// Get retrieves the latest snapshot for a profile.
func (s *SnapshotStore) Get(profile string) (Snapshot, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	snap, ok := s.data[profile]
	return snap, ok
}

// Delete removes the snapshot for a profile and flushes to disk.
func (s *SnapshotStore) Delete(profile string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, profile)
	return s.flush()
}

func (s *SnapshotStore) load() error {
	b, err := os.ReadFile(s.path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("snapshot: read %s: %w", s.path, err)
	}
	return json.Unmarshal(b, &s.data)
}

func (s *SnapshotStore) flush() error {
	b, err := json.MarshalIndent(s.data, "", "  ")
	if err != nil {
		return fmt.Errorf("snapshot: marshal: %w", err)
	}
	return os.WriteFile(s.path, b, 0600)
}
