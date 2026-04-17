package trend

import (
	"context"

	"github.com/user/portwatch/internal/scanner"
)

// Event carries a snapshot of the trend after each scan cycle.
type Event struct {
	OpenCount int
	Dir       Direction
}

// NewPipeline feeds port counts from in into the Tracker and forwards
// trend Events to the returned channel. The channel is closed when ctx
// is cancelled or in is closed.
func NewPipeline(ctx context.Context, tr *Tracker, in <-chan []scanner.Port) <-chan Event {
	out := make(chan Event, 8)
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
				tr.Record(len(ports))
				ev := Event{
					OpenCount: len(ports),
					Dir:       tr.Direction(),
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
