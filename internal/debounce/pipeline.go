package debounce

import (
	"context"
	"time"

	"portwatch/internal/scanner"
)

// PortChange represents an open/close transition detected by the monitor.
type PortChange struct {
	Port   scanner.Port
	Opened bool
}

// Pipeline wraps a Debouncer and bridges raw PortChange channels into
// debounced Event output, running until ctx is cancelled.
type Pipeline struct {
	d   *Debouncer
	in  <-chan PortChange
}

// NewPipeline creates a Pipeline with the given debounce delay and input channel.
func NewPipeline(delay time.Duration, in <-chan PortChange) *Pipeline {
	return &Pipeline{
		d:  New(delay),
		in: in,
	}
}

// C returns the debounced output channel.
func (p *Pipeline) C() <-chan Event {
	return p.d.C()
}

// Run reads from the input channel and pushes events into the debouncer
// until ctx is done or the input channel is closed.
func (p *Pipeline) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case ch, ok := <-p.in:
			if !ok {
				return
			}
			p.d.Push(Event{
				Port:   ch.Port,
				Opened: ch.Opened,
			})
		}
	}
}

// RunAndClose runs the pipeline and closes the debouncer when finished.
// This is a convenience wrapper around Run for callers that own the
// pipeline lifecycle and want the output channel closed on exit.
func (p *Pipeline) RunAndClose(ctx context.Context) {
	p.Run(ctx)
	p.d.Close()
}
