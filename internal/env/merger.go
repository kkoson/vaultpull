// Package env provides utilities for reading and writing .env files.
package env

import (
	"fmt"
	"os"
	"strings"
)

// MergeMode controls how existing keys are handled during a merge.
type MergeMode int

const (
	// MergeOverwrite replaces existing keys with new values.
	MergeOverwrite MergeMode = iota
	// MergeKeepExisting preserves existing keys and only adds new ones.
	MergeKeepExisting
)

// Merger merges a set of secrets into an existing .env file,
// preserving comments and unmanaged keys.
type Merger struct {
	reader *Reader
	writer *Writer
	mode   MergeMode
}

// NewMerger creates a Merger for the given file path and merge mode.
func NewMerger(path string, mode MergeMode) *Merger {
	return &Merger{
		reader: NewReader(path),
		writer: NewWriter(path),
		mode:   mode,
	}
}

// Merge reads the existing .env file, applies incoming secrets according
// to the MergeMode, and writes the result back to disk.
// If the file does not exist it is created with only the incoming secrets.
func (m *Merger) Merge(incoming map[string]string) error {
	existing, err := m.reader.Read()
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("merger: read existing file: %w", err)
	}

	merged := make(map[string]string, len(existing)+len(incoming))

	// Seed with existing values.
	for k, v := range existing {
		merged[k] = v
	}

	// Apply incoming according to mode.
	for k, v := range incoming {
		normKey := strings.ToUpper(k)
		if m.mode == MergeKeepExisting {
			if _, exists := merged[normKey]; exists {
				continue
			}
		}
		merged[normKey] = v
	}

	if err := m.writer.Write(merged); err != nil {
		return fmt.Errorf("merger: write merged file: %w", err)
	}
	return nil
}
