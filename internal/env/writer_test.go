package env

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func tmpFile(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), ".env")
}

func TestWrite_BasicSecrets(t *testing.T) {
	p := tmpFile(t)
	w := NewWriter(p)

	err := w.Write(map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(p)
	content := string(data)
	if !strings.Contains(content, "DB_HOST=localhost") {
		t.Errorf("expected DB_HOST=localhost in output, got:\n%s", content)
	}
	if !strings.Contains(content, "DB_PORT=5432") {
		t.Errorf("expected DB_PORT=5432 in output, got:\n%s", content)
	}
}

func TestWrite_QuotesValueWithSpaces(t *testing.T) {
	p := tmpFile(t)
	w := NewWriter(p)

	err := w.Write(map[string]string{"MSG": "hello world"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(p)
	if !strings.Contains(string(data), `MSG="hello world"`) {
		t.Errorf("expected quoted value, got: %s", data)
	}
}

func TestWrite_EmptyMap(t *testing.T) {
	p := tmpFile(t)
	w := NewWriter(p)

	if err := w.Write(map[string]string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(p)
	if len(data) != 0 {
		t.Errorf("expected empty file, got: %s", data)
	}
}

func TestQuoteValue(t *testing.T) {
	cases := []struct {
		input, want string
	}{
		{"simple", "simple"},
		{"with space", `"with space"`},
		{"has#hash", `"has#hash"`},
		{`has"quote`, `"has\"quote"`},
	}
	for _, c := range cases {
		got := quoteValue(c.input)
		if got != c.want {
			t.Errorf("quoteValue(%q) = %q, want %q", c.input, got, c.want)
		}
	}
}
