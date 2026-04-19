package portdiff

import (
	"context"

	"portwatch/internal/scanner"
)

// Event carries a diff produced from consecutive scanner samples.
type Event struct {
	Summary Summary
}

// NewPipeline reads successive port lists from in, computes diffs, and emits
// non-empty summaries on the returned channel. The channel is closed when ctx
// is cancelled or in is closed.
func NewPipeline(ctx context.Context, in <-chan []scanner.Port) <-chan Event {
	out := make(chan Event)
	go func() {
		defer close(out)
		var prev []scanner.Port
		for {
			select {
			case <-ctx.Done():
				return
			case ports, ok := <-in:
				if !ok {
					return
				}
				if prev != nil {
					s := Compare(prev, ports)
					if !s.Empty() {
						select {
						case out <- Event{Summary: s}:
						case <-ctx.Done():
							return
						}
					}
				}
				prev = ports
			}
		}
	}()
	return out
}
