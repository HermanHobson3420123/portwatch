package ratelimit

import (
	"context"

	"portwatch/internal/monitor"
)

// Pipeline wraps a Limiter and forwards only non-suppressed events.
type Pipeline struct {
	limiter *Limiter
	in      <-chan monitor.Event
	out     chan monitor.Event
}

// NewPipeline creates a Pipeline that rate-limits events from in.
func NewPipeline(l *Limiter, in <-chan monitor.Event) *Pipeline {
	return &Pipeline{
		limiter: l,
		in:      in,
		out:     make(chan monitor.Event, 64),
	}
}

// Out returns the filtered output channel.
func (p *Pipeline) Out() <-chan monitor.Event {
	return p.out
}

// Run reads from the input channel, suppresses events that exceed the rate
// limit, and forwards the rest. It exits when ctx is cancelled or in closes.
func (p *Pipeline) Run(ctx context.Context) {
	defer close(p.out)
	for {
		select {
		case <-ctx.Done():
			return
		case ev, ok := <-p.in:
			if !ok {
				return
			}
			key := ev.Port.String()
			if p.limiter.Allow(key) {
				select {
				case p.out <- ev:
				case <-ctx.Done():
					return
				}
			}
		}
	}
}
