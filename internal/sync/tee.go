package sync

import "errors"

// SecretWriter is the interface satisfied by env.Writer.
type SecretWriter interface {
	Write(secrets map[string]string) error
}

// TeeWriter fans out a Write call to multiple SecretWriter targets.
type TeeWriter struct {
	writers []SecretWriter
}

// NewTeeWriter creates a TeeWriter that writes to all provided writers.
// Panics if writers is empty.
func NewTeeWriter(w ...SecretWriter) *TeeWriter {
	if len(w) == 0 {
		panic("tee: at least one writer is required")
	}
	return &TeeWriter{writers: w}
}

// Write calls Write on every underlying writer.
// Returns the first error encountered; remaining writers are still attempted.
func (t *TeeWriter) Write(secrets map[string]string) error {
	var firstErr error
	for _, w := range t.writers {
		if err := w.Write(secrets); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

// Len returns the number of underlying writers.
func (t *TeeWriter) Len() int { return len(t.writers) }

// ErrTeeEmpty is returned when NewTeeWriter is called with no writers.
var ErrTeeEmpty = errors.New("tee: no writers provided")
