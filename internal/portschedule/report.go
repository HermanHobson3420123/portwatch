package portschedule

import (
	"fmt"
	"io"
	"time"
)

// PrintSchedule writes a human-readable table of all schedule entries to w.
func PrintSchedule(w io.Writer, s *Schedule) {
	entries := s.All()
	if len(entries) == 0 {
		fmt.Fprintln(w, "no scheduled windows defined")
		return
	}
	now := time.Now()
	fmt.Fprintf(w, "%-20s %-8s %6s-%6s %-8s %s\n", "LABEL", "PROTO", "LOW", "HIGH", "STATUS", "WINDOW")
	for _, e := range entries {
		status := "inactive"
		if e.Active(now) {
			status = "active"
		}
		fmt.Fprintf(w, "%-20s %-8s %6d-%6d %-8s %s -> %s\n",
			e.Label, e.Protocol, e.PortLow, e.PortHigh, status,
			e.Start.Format(time.RFC3339), e.End.Format(time.RFC3339))
	}
}

// Summary returns a one-line summary of active vs total windows.
func Summary(s *Schedule) string {
	all := s.All()
	active := s.ActiveNow(time.Now())
	return fmt.Sprintf("%d/%d schedule windows active", len(active), len(all))
}
