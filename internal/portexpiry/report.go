package portexpiry

import (
	"fmt"
	"io"
	"sort"
	"time"
)

// PrintExpired writes a human-readable summary of expired entries to w.
// Entries are sorted by port number before printing.
func PrintExpired(w io.Writer, entries []Entry, now time.Time) {
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Port.Number < entries[j].Port.Number
	})
	fmt.Fprintf(w, "%-8s %-6s %s\n", "PROTO", "PORT", "LAST SEEN")
	for _, e := range entries {
		age := now.Sub(e.LastSeen).Round(time.Second)
		fmt.Fprintf(w, "%-8s %-6d %s ago\n", e.Port.Proto, e.Port.Number, age)
	}
}

// Summary returns a one-line string describing the expiry state.
func Summary(entries []Entry) string {
	if len(entries) == 0 {
		return "no expired ports"
	}
	return fmt.Sprintf("%d port(s) expired", len(entries))
}

// PrintSummaryLine writes the summary line produced by Summary to w,
// followed by a newline. It returns any write error encountered.
func PrintSummaryLine(w io.Writer, entries []Entry) error {
	_, err := fmt.Fprintln(w, Summary(entries))
	return err
}
