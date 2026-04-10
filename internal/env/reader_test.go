package env

import (
	"os"
	"testing"
)

func writeEnvFile(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestRead_BasicPairs(t *testing.T) {
	path := writeEnvFile(t, "FOO=bar\nBAZ=qux\n")
	r := NewReader(path)
	got, err := r.Read()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["FOO"] != "bar" || got["BAZ"] != "qux" {
		t.Errorf("unexpected result: %v", got)
	}
}

func TestRead_IgnoresComments(t *testing.T) {
	path := writeEnvFile(t, "# comment\nKEY=value\n")
	r := NewReader(path)
	got, err := r.Read()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := got["# comment"]; ok {
		t.Error("comment line should not be parsed as key")
	}
	if got["KEY"] != "value" {
		t.Errorf("expected KEY=value, got %v", got)
	}
}

func TestRead_UnquotesValues(t *testing.T) {
	path := writeEnvFile(t, `SECRET="hello world"` + "\n")
	r := NewReader(path)
	got, err := r.Read()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["SECRET"] != "hello world" {
		t.Errorf("expected unquoted value, got %q", got["SECRET"])
	}
}

func TestRead_FileNotExist(t *testing.T) {
	r := NewReader("/nonexistent/path/.env")
	got, err := r.Read()
	if err != nil {
		t.Fatalf("expected no error for missing file, got: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected empty map for missing file, got %v", got)
	}
}

func TestRead_EmptyFile(t *testing.T) {
	path := writeEnvFile(t, "")
	r := NewReader(path)
	got, err := r.Read()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected empty map, got %v", got)
	}
}
