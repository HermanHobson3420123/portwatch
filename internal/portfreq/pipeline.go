package portfreq

import (
	"context"
	"time"

	"portwatch/internal/scanner"
)

// NewPipeline feeds every scan result into the Tracker and forwards the
// slice unchanged so downstream stages are unaffected.
func NewPipeline(ctx context.Context, tr *Tracker, in <-chan []scanner.Port) <-chan []scanner.Port {
	out := make(chan []scanner.Port)
	go func() {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				return
			case ports, ok := <-in:
				if !ok {
					return
				}
				now := time.Now()
				for _, p := range ports {
					tr.Record(p, now)
				}
				select {
				case out <- ports:
				case <-ctx.Done():
					return
				}
			}
		}
	}()
	return out
}
