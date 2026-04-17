package sampler

import (
	"context"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Sample holds a timestamped port snapshot.
type Sample struct {
	Time  time.Time
	Ports []scanner.Port
}

// Sampler periodically collects port snapshots.
type Sampler struct {
	scan     func(ctx context.Context) ([]scanner.Port, error)
	interval time.Duration
	out      chan Sample
}

// New creates a Sampler that calls scan every interval.
func New(scan func(ctx context.Context) ([]scanner.Port, error), interval time.Duration) *Sampler {
	return &Sampler{
		scan:     scan,
		interval: interval,
		out:      make(chan Sample, 8),
	}
}

// Out returns the read-only channel of samples.
func (s *Sampler) Out() <-chan Sample { return s.out }

// Run collects samples until ctx is cancelled.
func (s *Sampler) Run(ctx context.Context) {
	defer close(s.out)
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case t := <-ticker.C:
			ports, err := s.scan(ctx)
			if err != nil {
				continue
			}
			select {
			case s.out <- Sample{Time: t, Ports: ports}:
			case <-ctx.Done():
				return
			}
		}
	}
}
