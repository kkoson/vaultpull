package sync

import (
	"fmt"
	"io"
	"os"
	"sync/atomic"
)

// Progress tracks and reports sync progress across multiple profiles.
type Progress struct {
	out     io.Writer
	total   int
	done    atomic.Int32
	failed  atomic.Int32
	verbose bool
}

// NewProgress creates a Progress reporter writing to out.
// If out is nil, os.Stdout is used.
func NewProgress(out io.Writer, total int, verbose bool) *Progress {
	if out == nil {
		out = os.Stdout
	}
	return &Progress{
		out:     out,
		total:   total,
		verbose: verbose,
	}
}

// ProfileStarted logs the start of a profile sync when verbose is enabled.
func (p *Progress) ProfileStarted(profile string) {
	if p.verbose {
		fmt.Fprintf(p.out, "[%d/%d] syncing profile %q...\n",
			int(p.done.Load())+1, p.total, profile)
	}
}

// ProfileDone records a successful profile completion and logs a summary line.
func (p *Progress) ProfileDone(profile string, added, updated, removed int) {
	p.done.Add(1)
	if p.verbose {
		fmt.Fprintf(p.out, "  ✓ %s: +%d ~%d -%d\n", profile, added, updated, removed)
	}
}

// ProfileFailed records a failed profile and prints the error.
func (p *Progress) ProfileFailed(profile string, err error) {
	p.done.Add(1)
	p.failed.Add(1)
	fmt.Fprintf(p.out, "  ✗ %s: %v\n", profile, err)
}

// Summary prints the final summary line.
func (p *Progress) Summary() {
	failed := int(p.failed.Load())
	succeeded := int(p.done.Load()) - failed
	if failed == 0 {
		fmt.Fprintf(p.out, "sync complete: %d/%d profiles succeeded\n", succeeded, p.total)
	} else {
		fmt.Fprintf(p.out, "sync complete: %d/%d profiles succeeded, %d failed\n",
			succeeded, p.total, failed)
	}
}

// FailedCount returns the number of profiles that failed.
func (p *Progress) FailedCount() int {
	return int(p.failed.Load())
}
