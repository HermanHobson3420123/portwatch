package portlabel

import (
	"context"

	"portwatch/internal/scanner"
)

// Annotated wraps a scanned port with its label.
type Annotated struct {
	Port  scanner.Port
	Label string
}

// NewPipeline reads ports from in, attaches labels, and forwards to the
// returned channel. The output channel is closed when in is closed or ctx
// is cancelled.
func NewPipeline(ctx context.Context, in <-chan scanner.Port) <-chan Annotated {
	out := make(chan Annotated)
	go func() {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				return
			case p, ok := <-in:
				if !ok {
					return
				}
				a := Annotated{
					Port:  p,
					Label: Annotate(uint16(p.Port), p.Protocol),
				}
				select {
				case out <- a:
				case <-ctx.Done():
					return
				}
			}
		}
	}()
	return out
}
