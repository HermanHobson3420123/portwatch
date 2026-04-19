package portping

import (
	"context"

	"portwatch/internal/scanner"
)

// NewPipeline reads port slices from in, pings each port, and forwards
// results to the returned channel. The channel is closed when in is
// exhausted or ctx is cancelled.
func NewPipeline(ctx context.Context, p *Pinger, in <-chan []scanner.Port) <-chan []Result {
	out := make(chan []Result)
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
				results := p.PingAll(ctx, ports)
				select {
				case out <- results:
				case <-ctx.Done():
					return
				}
			}
		}
	}()
	return out
}
