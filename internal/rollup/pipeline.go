package rollup

import (
	"context"

	"portwatch/internal/monitor"
)

// Pipeline wires a monitor.Monitor diff channel into a Rollup, then
// forwards the resulting batched Events to a caller-supplied handler.
type Pipeline struct {
	rollup *Rollup
}

// NewPipeline creates a Pipeline backed by the given Rollup.
func NewPipeline(r *Rollup) *Pipeline {
	return &Pipeline{rollup: r}
}

// Run reads monitor.Diff values from diffs, feeds them into the Rollup,
// and calls handler for every batched Event.  It blocks until ctx is
// cancelled or diffs is closed.
func (p *Pipeline) Run(ctx context.Context, diffs <-chan monitor.Diff, handler func(Event)) {
	go p.rollup.Watch(ctx, handler)

	for {
		select {
		case <-ctx.Done():
			return
		case d, ok := <-diffs:
			if !ok {
				return
			}
			for _, port := range d.Opened {
				p.rollup.AddOpened(port)
			}
			for _, port := range d.Closed {
				p.rollup.AddClosed(port)
			}
		}
	}
}
