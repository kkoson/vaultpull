package sync

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func tmpCheckpointPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "checkpoint.json")
}

func TestCheckpoint_SetAndGet(t *testing.T) {
	cp, err := NewCheckpoint(tmpCheckpointPath(t))
	if err != nil {
		t.Fatalf("NewCheckpoint: %v", err)
	}
	entry := CheckpointEntry{
		Profile:    "prod",
		SyncedAt:   time.Now().UTC().Truncate(time.Second),
		SecretPath: "secret/prod",
	}
	if err := cp.Set(entry); err != nil {
		t.Fatalf("Set: %v", err)
	}
	got, ok := cp.Get("prod")
	if !ok {
		t.Fatal("expected entry to exist")
	}
	if got.Profile != entry.Profile || got.SecretPath != entry.SecretPath {
		t.Errorf("got %+v, want %+v", got, entry)
	}
}

func TestCheckpoint_PersistsAcrossInstances(t *testing.T) {
	path := tmpCheckpointPath(t)
	cp1, _ := NewCheckpoint(path)
	_ = cp1.Set(CheckpointEntry{Profile: "dev", SyncedAt: time.Now(), SecretPath: "secret/dev"})

	cp2, err := NewCheckpoint(path)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	_, ok := cp2.Get("dev")
	if !ok {
		t.Error("expected persisted entry to be loaded")
	}
}

func TestCheckpoint_Clear(t *testing.T) {
	path := tmpCheckpointPath(t)
	cp, _ := NewCheckpoint(path)
	_ = cp.Set(CheckpointEntry{Profile: "staging", SyncedAt: time.Now(), SecretPath: "secret/staging"})
	if err := cp.Clear("staging"); err != nil {
		t.Fatalf("Clear: %v", err)
	}
	_, ok := cp.Get("staging")
	if ok {
		t.Error("expected entry to be removed")
	}
}

func TestCheckpoint_MissingFile_NoError(t *testing.T) {
	_, err := NewCheckpoint(filepath.Join(t.TempDir(), "missing", "cp.json"))
	// Missing parent dir is created on flush; missing file itself is fine.
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestCheckpoint_Get_NotFound(t *testing.T) {
	cp, _ := NewCheckpoint(tmpCheckpointPath(t))
	_, ok := cp.Get("nonexistent")
	if ok {
		t.Error("expected false for unknown profile")
	}
}

func TestCheckpoint_FilePermissions(t *testing.T) {
	path := tmpCheckpointPath(t)
	cp, _ := NewCheckpoint(path)
	_ = cp.Set(CheckpointEntry{Profile: "p", SyncedAt: time.Now(), SecretPath: "s"})
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if perm := info.Mode().Perm(); perm != 0o600 {
		t.Errorf("expected 0600, got %o", perm)
	}
}
