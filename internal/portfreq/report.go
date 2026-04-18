package portfreq

import (
	"fmt"
	"io"
	"text/tabwriter"
)

// PrintTop writes the top n entries from the tracker to w in a tabular format.
func PrintTop(w io.Writer, tr *Tracker, n int) {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "PROTO\tADDR\tCOUNT\tFIRST SEEN\tLAST SEEN")
	for _, e := range tr.Top(n) {
		fmt.Fprintf(tw, "%s\t%s\t%d\t%s\t%s\n",
			e.Port.Proto,
			e.Port.Addr,
			e.SeenCount,
			e.FirstSeen.Format("2006-01-02 15:04:05"),
			e.LastSeen.Format("2006-01-02 15:04:05"),
		)
	}
	tw.Flush()
}

// Summary returns a one-line string describing the most-seen port.
func Summary(tr *Tracker) string {
	top := tr.Top(1)
	if len(top) == 0 {
		return "no frequency data"
	}
	e := top[0]
	return fmt.Sprintf("most frequent: %s %s (%d scans)", e.Port.Proto, e.Port.Addr, e.SeenCount)
}
