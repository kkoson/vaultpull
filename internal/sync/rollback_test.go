package sync

import (
	"os"
	"path/filepath"
	"testing"
)

func tmpRollbackDir(t *testing.T) string {
	t.Helper()
	return t.TempDir()
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("writeFile: %v", err)
	}
}

func TestRollbackStore_BackupAndRestore(t *testing.T) {
	dir := tmpRollbackDir(t)
	store, err := NewRollbackStore(dir)
	if err != nil {
		t.Fatalf("NewRollbackStore: %v", err)
	}

	src := filepath.Join(dir, "dev.env")
	writeFile(t, src, "KEY=value\n")

	if err := store.Backup("dev", src); err != nil {
		t.Fatalf("Backup: %v", err)
	}
	if !store.HasBackup("dev") {
		t.Fatal("expected HasBackup to return true")
	}

	// Overwrite src then restore.
	writeFile(t, src, "KEY=changed\n")
	if err := store.Restore("dev", src); err != nil {
		t.Fatalf("Restore: %v", err)
	}

	got, _ := os.ReadFile(src)
	if string(got) != "KEY=value\n" {
		t.Errorf("restored content = %q; want %q", string(got), "KEY=value\n")
	}
}

func TestRollbackStore_Backup_SourceMissing_NoError(t *testing.T) {
	store, _ := NewRollbackStore(t.TempDir())
	// Source file does not exist — should be a no-op.
	if err := store.Backup("missing", "/nonexistent/path.env"); err != nil {
		t.Fatalf("expected no error for missing source, got: %v", err)
	}
	if store.HasBackup("missing") {
		t.Fatal("expected no backup to be created")
	}
}

func TestRollbackStore_Restore_NoBackup_Error(t *testing.T) {
	store, _ := NewRollbackStore(t.TempDir())
	if err := store.Restore("ghost", "/tmp/ghost.env"); err == nil {
		t.Fatal("expected error when no backup exists")
	}
}

func TestRollbackStore_Clear(t *testing.T) {
	dir := t.TempDir()
	store, _ := NewRollbackStore(dir)

	src := filepath.Join(dir, "prod.env")
	writeFile(t, src, "X=1\n")
	_ = store.Backup("prod", src)

	if err := store.Clear("prod"); err != nil {
		t.Fatalf("Clear: %v", err)
	}
	if store.HasBackup("prod") {
		t.Fatal("expected backup to be removed")
	}
}

func TestRollbackStore_Clear_NoBackup_NoError(t *testing.T) {
	// Clearing a key that was never backed up should not return an error.
	store, _ := NewRollbackStore(t.TempDir())
	if err := store.Clear("nonexistent"); err != nil {
		t.Fatalf("expected no error clearing non-existent backup, got: %v", err)
	}
}

func TestRollbackStore_BackupTimestamp_NoBackup(t *testing.T) {
	store, _ := NewRollbackStore(t.TempDir())
	ts := store.BackupTimestamp("nope")
	if !ts.IsZero() {
		t.Errorf("expected zero time, got %v", ts)
	}
}

func TestRollbackStore_BackupTimestamp_WithBackup(t *testing.T) {
	dir := t.TempDir()
	store, _ := NewRollbackStore(dir)

	src := filepath.Join(dir, "stg.env")
	writeFile(t, src, "A=b\n")
	_ = store.Backup("stg", src)

	ts := store.BackupTimestamp("stg")
	if ts.IsZero() {
		t.Error("expected non-zero timestamp after backup")
	}
}
