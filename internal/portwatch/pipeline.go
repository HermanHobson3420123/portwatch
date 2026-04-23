package portwatch

import (
	"context"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// PipelineConfig holds options for NewPipeline.
type PipelineConfig struct {
	Interval time.Duration
}

// DefaultPipelineConfig returns sensible defaults.
func DefaultPipelineConfig() PipelineConfig {
	return PipelineConfig{
		Interval: 5 * time.Second,
	}
}

// NewPipeline wires a Scanner into a Watcher and returns the event channel.
// The caller is responsible for consuming the channel until it is closed.
func NewPipeline(ctx context.Context, s *scanner.Scanner, cfg PipelineConfig) <-chan WatchEvent {
	w := New(s, cfg.Interval)
	return w.Watch(ctx)
}
