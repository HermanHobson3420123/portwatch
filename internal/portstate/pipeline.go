package portstate

import (
	"context"
	"time"

	"portwatch/internal/monitor"
)

// NewPipeline feeds monitor diff events into a Tracker and forwards them
// unchanged so downstream consumers can still process the same events.
func NewPipeline(ctx context.Context, tr *Tracker, in <-chan monitor.Event) <-chan monitor.Event {
	out := make(chan monitor.Event)
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
				now := time.Now()
				switch ev.Type {
				case monitor.Opened:
					tr.Open(ev.Port, now)
				case monitor.Closed:
					tr.Close(ev.Port)
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
