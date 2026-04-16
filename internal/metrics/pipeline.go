package metrics

import (
	"context"

	"github.com/user/portwatch/internal/scanner"
)

// Event carries the direction of a port change.
type Event struct {
	Opened bool
	Port   scanner.Port
}

// Pipeline wires a Collector to an incoming event channel, recording
// each opened/closed alert and forwarding events to the out channel.
func NewPipeline(ctx context.Context, c *Collector, in <-chan Event) <-chan Event {
	out := make(chan Event, 64)
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
				if ev.Opened {
					c.RecordOpened()
				} else {
					c.RecordClosed()
				}
				select {
				case out <- ev:
				case <-ctx.Done():
					return
				}
			}
		}
	}()
	return out
}
