package watchdog

import (
	"context"

	"portwatch/internal/scanner"
)

// ScanFunc is the signature of a function that performs a single scan cycle
// and returns the discovered ports or an error.
type ScanFunc func(ctx context.Context) ([]scanner.Port, error)

// Pipeline wraps a ScanFunc, feeds results into an output channel, and calls
// Beat on the Watchdog after every successful scan.
type Pipeline struct {
	scan     ScanFunc
	watchdog *Watchdog
	out      chan<- []scanner.Port
}

// NewPipeline creates a Pipeline that will send scan results to out.
func NewPipeline(scan ScanFunc, wd *Watchdog, out chan<- []scanner.Port) *Pipeline {
	return &Pipeline{
		scan:     scan,
		watchdog: wd,
		out:      out,
	}
}

// Run executes scan in a loop until ctx is cancelled. Each successful scan
// result is forwarded to the output channel and the watchdog is notified.
func (p *Pipeline) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		ports, err := p.scan(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			continue
		}

		p.watchdog.Beat()

		select {
		case p.out <- ports:
		case <-ctx.Done():
			return
		}
	}
}
