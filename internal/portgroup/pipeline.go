package portgroup

import (
	"context"

	"portwatch/internal/scanner"
)

// Sample carries a snapshot of ports together with their computed groups.
type Sample struct {
	Ports  []scanner.Port
	Groups []Group
}

// NewPipeline reads port slices from in, groups them with g, and forwards
// Sample values to the returned channel. The channel is closed when ctx is
// cancelled or in is closed.
func NewPipeline(ctx context.Context, g *Grouper, in <-chan []scanner.Port) <-chan Sample {
	out := make(chan Sample)
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
				s := Sample{
					Ports:  ports,
					Groups: g.Group(ports),
				}
				select {
				case out <- s:
				case <-ctx.Done():
					return
				}
			}
		}
	}()
	return out
}
