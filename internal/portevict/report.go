package portevict

import (
	"fmt"
	"io"
	"strings"
)

// PrintEvictions writes a human-readable eviction report to w.
func PrintEvictions(w io.Writer, entries []Entry) {
	if len(entries) == 0 {
		fmt.Fprintln(w, "no evictions recorded")
		return
	}
	fmt.Fprintf(w, "%-8s %-6s %-20s %-20s %s\n", "PORT", "PROTO", "OPENED", "CLOSED", "DURATION")
	fmt.Fprintln(w, strings.Repeat("-", 72))
	for _, e := range entries {
		fmt.Fprintf(w, "%-8d %-6s %-20s %-20s %s\n",
			e.Port.Number,
			e.Port.Protocol,
			e.OpenedAt.Format("15:04:05.000"),
			e.ClosedAt.Format("15:04:05.000"),
			e.Duration.Round(1*1000*1000).String(),
		)
	}
}

// Summary returns a one-line summary string.
func Summary(entries []Entry) string {
	if len(entries) == 0 {
		return "evictions: none"
	}
	return fmt.Sprintf("evictions: %d short-lived port(s) detected", len(entries))
}
