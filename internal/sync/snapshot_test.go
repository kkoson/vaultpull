package sync

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func tmpSnapshotPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "snapshots.json")
}

func TestSnapshotStore_SaveAndGet(t *testing.T) {
	store, err := NewSnapshotStore(tmpSnapshotPath(t))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	snap := Snapshot{
		Profile:    "prod",
		CapturedAt: time.Now().UTC(),
		Secrets:    map[string]string{"DB_PASS": "secret"},
	}
	if err := store.Save(snap); err != nil {
		t.Fatalf("Save: %v", err)
	}
	got, ok := store.Get("prod")
	if !ok {
		t.Fatal("expected snapshot to exist")
	}
	if got.Secrets["DB_PASS"] != "secret" {
		t.Errorf("expected DB_PASS=secret, got %s", got.Secrets["DB_PASS"])
	}
}

func TestSnapshotStore_PersistsAcrossInstances(t *testing.T) {
	path := tmpSnapshotPath(t)
	store1, _ := NewSnapshotStore(path)
	_ = store1.Save(Snapshot{Profile: "staging", Secrets: map[string]string{"KEY": "val"}})

	store2, err := NewSnapshotStore(path)
	if err != nil {
		t.Fatalf("reload error: %v", err)
	}
	snap, ok := store2.Get("staging")
	if !ok {
		t.Fatal("expected snapshot after reload")
	}
	if snap.Secrets["KEY"] != "val" {
		t.Errorf("unexpected value: %s", snap.Secrets["KEY"])
	}
}

func TestSnapshotStore_Delete(t *testing.T) {
	path := tmpSnapshotPath(t)
	store, _ := NewSnapshotStore(path)
	_ = store.Save(Snapshot{Profile: "dev", Secrets: map[string]string{"X": "1"}})
	if err := store.Delete("dev"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	_, ok := store.Get("dev")
	if ok {
		t.Error("expected snapshot to be deleted")
	}
}

func TestSnapshotStore_MissingFile_NoError(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nonexistent.json")
	store, err := NewSnapshotStore(path)
	if err != nil {
		t.Fatalf("expected no error for missing file, got: %v", err)
	}
	_, ok := store.Get("any")
	if ok {
		t.Error("expected empty store")
	}
}

func TestSnapshotStore_InvalidFile_ReturnsError(t *testing.T) {
	path := tmpSnapshotPath(t)
	_ = os.WriteFile(path, []byte("not-json"), 0600)
	_, err := NewSnapshotStore(path)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}
