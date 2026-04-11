// Package audit provides structured audit logging for sync operations.
package audit

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

// Entry represents a single audit log entry.
type Entry struct {
	Timestamp  time.Time `json:"timestamp"`
	Profile    string    `json:"profile"`
	VaultPath  string    `json:"vault_path"`
	OutputFile string    `json:"output_file"`
	Added      int       `json:"added"`
	Updated    int       `json:"updated"`
	Removed    int       `json:"removed"`
	Unchanged  int       `json:"unchanged"`
	DryRun     bool      `json:"dry_run"`
	Error      string    `json:"error,omitempty"`
}

// Logger writes audit entries to a destination.
type Logger struct {
	w io.Writer
}

// NewLogger creates a Logger writing to the given writer.
// Pass nil to use os.Stdout.
func NewLogger(w io.Writer) *Logger {
	if w == nil {
		w = os.Stdout
	}
	return &Logger{w: w}
}

// Write serialises the entry as JSON and writes it followed by a newline.
func (l *Logger) Write(e Entry) error {
	if e.Timestamp.IsZero() {
		e.Timestamp = time.Now().UTC()
	}
	b, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("audit: marshal entry: %w", err)
	}
	_, err = fmt.Fprintf(l.w, "%s\n", b)
	return err
}
