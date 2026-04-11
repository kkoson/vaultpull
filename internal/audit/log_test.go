package audit_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/yourusername/vaultpull/internal/audit"
)

func TestWrite_BasicEntry(t *testing.T) {
	var buf bytes.Buffer
	l := audit.NewLogger(&buf)

	e := audit.Entry{
		Timestamp:  time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		Profile:    "staging",
		VaultPath:  "secret/app",
		OutputFile: ".env",
		Added:      3,
		Updated:    1,
		Removed:    0,
		Unchanged:  5,
		DryRun:     false,
	}

	if err := l.Write(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	line := strings.TrimSpace(buf.String())
	var got audit.Entry
	if err := json.Unmarshal([]byte(line), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if got.Profile != "staging" {
		t.Errorf("profile: want staging, got %s", got.Profile)
	}
	if got.Added != 3 {
		t.Errorf("added: want 3, got %d", got.Added)
	}
}

func TestWrite_SetsTimestampWhenZero(t *testing.T) {
	var buf bytes.Buffer
	l := audit.NewLogger(&buf)

	e := audit.Entry{Profile: "prod"}
	if err := l.Write(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var got audit.Entry
	if err := json.Unmarshal(bytes.TrimSpace(buf.Bytes()), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestWrite_WithError(t *testing.T) {
	var buf bytes.Buffer
	l := audit.NewLogger(&buf)

	e := audit.Entry{Profile: "dev", Error: "vault unreachable"}
	if err := l.Write(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var got audit.Entry
	if err := json.Unmarshal(bytes.TrimSpace(buf.Bytes()), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.Error != "vault unreachable" {
		t.Errorf("error field: want 'vault unreachable', got %q", got.Error)
	}
}

func TestNewLogger_NilUsesStdout(t *testing.T) {
	// Should not panic
	l := audit.NewLogger(nil)
	if l == nil {
		t.Fatal("expected non-nil logger")
	}
}
