// Package portannounce emits a human-readable startup summary of currently
// open ports so operators know the initial state when portwatch begins.
package portannounce

import (
	"fmt"
	"io"
	"sort"
	"strings"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Summary holds the announced state.
type Summary struct {
	At    time.Time
	Ports []scanner.Port
}

// Announcer writes a startup port summary to a writer.
type Announcer struct {
	out io.Writer
}

// New returns an Announcer that writes to w.
func New(w io.Writer) *Announcer {
	return &Announcer{out: w}
}

// Announce writes a formatted summary of ports to the configured writer.
func (a *Announcer) Announce(ports []scanner.Port) Summary {
	now := time.Now().UTC()
	sorted := make([]scanner.Port, len(ports))
	copy(sorted, ports)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].Protocol != sorted[j].Protocol {
			return sorted[i].Protocol < sorted[j].Protocol
		}
		return sorted[i].Number < sorted[j].Number
	})

	fmt.Fprintf(a.out, "portwatch startup — %s\n", now.Format(time.RFC3339))
	fmt.Fprintf(a.out, "%d port(s) currently open:\n", len(sorted))
	for _, p := range sorted {
		fmt.Fprintf(a.out, "  %-5s %d\n", strings.ToUpper(p.Protocol), p.Number)
	}
	if len(sorted) == 0 {
		fmt.Fprintln(a.out, "  (none)")
	}

	return Summary{At: now, Ports: sorted}
}
