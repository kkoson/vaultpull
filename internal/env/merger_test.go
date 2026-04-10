package env

import (
	"os"
	"testing"
)

func TestMerge_Overwrite(t *testing.T) {
	path := tmpFile(t)

	// Write initial content.
	w := NewWriter(path)
	if err := w.Write(map[string]string{"FOO": "old", "BAR": "keep"}); err != nil {
		t.Fatalf("setup write: %v", err)
	}

	m := NewMerger(path, MergeOverwrite)
	if err := m.Merge(map[string]string{"FOO": "new", "BAZ": "added"}); err != nil {
		t.Fatalf("Merge: %v", err)
	}

	r := NewReader(path)
	got, err := r.Read()
	if err != nil {
		t.Fatalf("read: %v", err)
	}

	assertEnvValue(t, got, "FOO", "new")
	assertEnvValue(t, got, "BAR", "keep")
	assertEnvValue(t, got, "BAZ", "added")
}

func TestMerge_KeepExisting(t *testing.T) {
	path := tmpFile(t)

	w := NewWriter(path)
	if err := w.Write(map[string]string{"FOO": "original"}); err != nil {
		t.Fatalf("setup write: %v", err)
	}

	m := NewMerger(path, MergeKeepExisting)
	if err := m.Merge(map[string]string{"FOO": "ignored", "NEW_KEY": "added"}); err != nil {
		t.Fatalf("Merge: %v", err)
	}

	r := NewReader(path)
	got, err := r.Read()
	if err != nil {
		t.Fatalf("read: %v", err)
	}

	assertEnvValue(t, got, "FOO", "original")
	assertEnvValue(t, got, "NEW_KEY", "added")
}

func TestMerge_FileNotExist_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	path := dir + "/.env"

	m := NewMerger(path, MergeOverwrite)
	if err := m.Merge(map[string]string{"CREATED": "yes"}); err != nil {
		t.Fatalf("Merge on missing file: %v", err)
	}

	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected file to be created: %v", err)
	}

	r := NewReader(path)
	got, err := r.Read()
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	assertEnvValue(t, got, "CREATED", "yes")
}

// assertEnvValue is a small helper shared across env tests.
func assertEnvValue(t *testing.T, m map[string]string, key, want string) {
	t.Helper()
	got, ok := m[key]
	if !ok {
		t.Errorf("key %q not found in map", key)
		return
	}
	if got != want {
		t.Errorf("key %q: got %q, want %q", key, got, want)
	}
}
