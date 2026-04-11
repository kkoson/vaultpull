package diff

import (
	"fmt"
	"io"
	"os"
	"sort"
)

// Color codes for terminal output.
const (
	colorReset  = "\033[0m"
	colorGreen  = "\033[32m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorGray   = "\033[90m"
)

// Printer writes a human-readable diff summary to an io.Writer.
type Printer struct {
	w     io.Writer
	color bool
}

// NewPrinter returns a Printer that writes to w.
// Color output is enabled when w is os.Stdout and the terminal supports it.
func NewPrinter(w io.Writer) *Printer {
	colorEnabled := w == os.Stdout
	return &Printer{w: w, color: colorEnabled}
}

// Print writes a formatted diff of the given Changes to the printer's writer.
func (p *Printer) Print(changes []Change) {
	if len(changes) == 0 {
		fmt.Fprintln(p.w, "No changes.")
		return
	}

	sort.Slice(changes, func(i, j int) bool {
		return changes[i].Key < changes[j].Key
	})

	for _, c := range changes {
		switch c.Type {
		case Added:
			fmt.Fprintf(p.w, "%s+ %s=%s%s\n", p.code(colorGreen), c.Key, c.NewValue, p.code(colorReset))
		case Removed:
			fmt.Fprintf(p.w, "%s- %s=%s%s\n", p.code(colorRed), c.Key, c.OldValue, p.code(colorReset))
		case Updated:
			fmt.Fprintf(p.w, "%s~ %s: %s → %s%s\n", p.code(colorYellow), c.Key, c.OldValue, c.NewValue, p.code(colorReset))
		case Unchanged:
			fmt.Fprintf(p.w, "%s  %s=%s%s\n", p.code(colorGray), c.Key, c.NewValue, p.code(colorReset))
		}
	}
}

// Summary prints a short count summary of changes.
func (p *Printer) Summary(changes []Change) {
	var added, removed, updated, unchanged int
	for _, c := range changes {
		switch c.Type {
		case Added:
			added++
		case Removed:
			removed++
		case Updated:
			updated++
		case Unchanged:
			unchanged++
		}
	}
	fmt.Fprintf(p.w, "Summary: %s+%d added%s, %s-%d removed%s, %s~%d updated%s, %d unchanged\n",
		p.code(colorGreen), added, p.code(colorReset),
		p.code(colorRed), removed, p.code(colorReset),
		p.code(colorYellow), updated, p.code(colorReset),
		unchanged,
	)
}

func (p *Printer) code(c string) string {
	if p.color {
		return c
	}
	return ""
}
