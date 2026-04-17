package sampler

import (
	"context"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// NewPipeline wires a Sampler using the provided scan function and interval,
// starts it in a goroutine, and returns the output channel.
func NewPipeline(
	ctx context.Context,
	scan func(ctx context.Context) ([]scanner.Port, error),
	interval time.Duration,
) <-chan Sample {
	s := New(scan, interval)
	go s.Run(ctx)
	return s.Out()
}
