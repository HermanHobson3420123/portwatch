package portburst

import (
	"context"

	"portwatch/internal/monitor"
)

// NewPipeline reads monitor.Event values from in, passes each port through
// the Detector, and forwards any burst Event to the returned channel.
// The output channel is closed when in is drained or ctx is cancelled.
func NewPipeline(ctx context.Context, d *Detector, in <-chan monitor.Event) <-chan Event {
	out := make(chan Event)
	go func() {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				return
			case ev, ok := <-in:
				if !ok {
					return
				}
				if burst := d.Record(ev.Port); burst != nil {
					select {
					case out <- *burst:
					case <-ctx.Done():
						return
					}
				}
			}
		}
	}()
	return out
}
