package portwindow

import (
	"fmt"
	"io"
	"sort"
	"strings"
	"time"
)

// PrintActive writes a human-readable summary of active window entries to w.
func PrintActive(out io.Writer, entries []Entry) {
	if len(entries) == 0 {
		fmt.Fprintln(out, "no ports active in window")
		return
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Port.Proto != entries[j].Port.Proto {
			return entries[i].Port.Proto < entries[j].Port.Proto
		}
		return entries[i].Port.Addr < entries[j].Port.Addr
	})
	fmt.Fprintf(out, "%-8s %-22s %-10s %-10s %s\n", "PROTO", "ADDR", "FIRST", "LAST", "COUNT")
	fmt.Fprintln(out, strings.Repeat("-", 64))
	for _, e := range entries {
		fmt.Fprintf(out, "%-8s %-22s %-10s %-10s %d\n",
			e.Port.Proto,
			e.Port.Addr,
			e.FirstSeen.Format(time.TimeOnly),
			e.LastSeen.Format(time.TimeOnly),
			e.Count,
		)
	}
}

// Summary returns a one-line summary string.
func Summary(entries []Entry) string {
	if len(entries) == 0 {
		return "window: 0 active ports"
	}
	total := 0
	for _, e := range entries {
		total += e.Count
	}
	return fmt.Sprintf("window: %d active port(s), %d total observations", len(entries), total)
}
