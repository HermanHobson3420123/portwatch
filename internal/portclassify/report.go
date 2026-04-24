package portclassify

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// Summary returns a one-line string summarising the tier distribution.
func Summary(classes []Class) string {
	counts := map[Tier]int{}
	for _, c := range classes {
		counts[c.Tier]++
	}
	return fmt.Sprintf("system=%d registered=%d dynamic=%d",
		counts[TierSystem], counts[TierRegistered], counts[TierDynamic])
}

// PrintReport writes a formatted classification table to w.
func PrintReport(w io.Writer, classes []Class) {
	if len(classes) == 0 {
		fmt.Fprintln(w, "no ports to classify")
		return
	}

	// Sort by port number for deterministic output.
	sorted := make([]Class, len(classes))
	copy(sorted, classes)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Port.Number < sorted[j].Port.Number
	})

	fmt.Fprintf(w, "%-8s %-6s %-12s %s\n", "PORT", "PROTO", "TIER", "SERVICE")
	fmt.Fprintln(w, strings.Repeat("-", 40))
	for _, c := range sorted {
		svc := c.Service
		if svc == "" {
			svc = "-"
		}
		fmt.Fprintf(w, "%-8d %-6s %-12s %s\n",
			c.Port.Number, c.Port.Protocol, c.Tier, svc)
	}
}
