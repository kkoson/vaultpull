package sync

import (
	"os"
	"path/filepath"
	"testing"
)

func tmpReplayDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "vaultpull-replay-*")
	if err != nil {
		t.Fatalf("tmpReplayDir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func TestReplayStore_SaveAndLoad(t *testing.T) {
	store := NewReplayStore(tmpReplayDir(t))
	secrets := map[string]string{"DB_PASS": "s3cr3t", "API_KEY": "abc123"}

	if err := store.Save("prod", secrets); err != nil {
		t.Fatalf("Save: %v", err)
	}

	got, ok := store.Load("prod")
	if !ok {
		t.Fatal("expected Load to return true")
	}
	if got["DB_PASS"] != "s3cr3t" {
		t.Errorf("DB_PASS: got %q", got["DB_PASS"])
	}
}

func TestReplayStore_Load_MissingFile(t *testing.T) {
	store := NewReplayStore(tmpReplayDir(t))
	_, ok := store.Load("nonexistent")
	if ok {
		t.Fatal("expected false for missing profile")
	}
}

func TestReplayStore_PersistsAcrossInstances(t *testing.T) {
	dir := tmpReplayDir(t)
	NewReplayStore(dir).Save("staging", map[string]string{"X": "1"})

	got, ok := NewReplayStore(dir).Load("staging")
	if !ok || got["X"] != "1" {
		t.Fatalf("expected X=1, got %v %v", got, ok)
	}
}

func TestReplayStore_Delete(t *testing.T) {
	dir := tmpReplayDir(t)
	store := NewReplayStore(dir)
	store.Save("dev", map[string]string{"K": "v"})

	if err := store.Delete("dev"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	_, ok := store.Load("dev")
	if ok {
		t.Fatal("expected false after delete")
	}
}

func TestReplayStore_Delete_NoFile_NoError(t *testing.T) {
	store := NewReplayStore(tmpReplayDir(t))
	if err := store.Delete("ghost"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestReplayStore_FileStoredUnderDir(t *testing.T) {
	dir := tmpReplayDir(t)
	store := NewReplayStore(dir)
	store.Save("alpha", map[string]string{})

	expected := filepath.Join(dir, "alpha.replay.json")
	if _, err := os.Stat(expected); err != nil {
		t.Fatalf("expected file %s: %v", expected, err)
	}
}
