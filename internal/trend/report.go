package trend

import (
	"fmt"
	"io"
	"text/tabwriter"
	"time"
)

// PrintSamples writes a formatted table of samples to w.
func PrintSamples(w io.Writer, samples []Sample) {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "TIME\tOPEN PORTS")
	for _, s := range samples {
		fmt.Fprintf(tw, "%s\t%d\n", s.At.Format(time.RFC3339), s.Count)
	}
	tw.Flush()
}

// Summary returns a one-line human-readable summary of the current trend.
func Summary(tr *Tracker) string {
	samples := tr.Samples()
	if len(samples) == 0 {
		return "no data"
	}
	last := samples[len(samples)-1]
	dir := tr.Direction()
	return fmt.Sprintf("open ports: %d  trend: %s  (samples: %d)",
		last.Count, dir, len(samples))
}
