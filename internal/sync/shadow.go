package sync

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// ShadowMode controls whether shadow writes are active.
type ShadowMode int

const (
	// ShadowOff disables shadow writing.
	ShadowOff ShadowMode = iota
	// ShadowOn enables shadow writing to a secondary destination.
	ShadowOn
)

// ShadowWriter writes secrets to a secondary "shadow" destination in parallel
// with the primary write, without affecting the primary result. It is useful
// for validating a new target before fully migrating to it.
type ShadowWriter struct {
	mu      sync.Mutex
	mode    ShadowMode
	dest    io.Writer
	records []ShadowRecord
}

// ShadowRecord captures a single shadow write attempt.
type ShadowRecord struct {
	Profile   string
	Timestamp time.Time
	Err       error
}

// NewShadowWriter creates a ShadowWriter. If dest is nil, os.Stderr is used.
func NewShadowWriter(mode ShadowMode, dest io.Writer) *ShadowWriter {
	if dest == nil {
		dest = os.Stderr
	}
	return &ShadowWriter{mode: mode, dest: dest}
}

// Write performs a shadow write for the given profile using fn.
// Errors from fn are recorded but never returned to the caller.
func (s *ShadowWriter) Write(profile string, fn func() error) {
	if s.mode == ShadowOff {
		return
	}
	var err error
	if fn != nil {
		err = fn()
	}
	rec := ShadowRecord{
		Profile:   profile,
		Timestamp: time.Now().UTC(),
		Err:       err,
	}
	s.mu.Lock()
	s.records = append(s.records, rec)
	s.mu.Unlock()
	if err != nil {
		fmt.Fprintf(s.dest, "[shadow] profile=%s error=%v\n", profile, err)
	}
}

// Records returns a snapshot of all recorded shadow attempts.
func (s *ShadowWriter) Records() []ShadowRecord {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]ShadowRecord, len(s.records))
	copy(out, s.records)
	return out
}

// ErrorCount returns the number of shadow writes that produced an error.
func (s *ShadowWriter) ErrorCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	count := 0
	for _, r := range s.records {
		if r.Err != nil {
			count++
		}
	}
	return count
}
