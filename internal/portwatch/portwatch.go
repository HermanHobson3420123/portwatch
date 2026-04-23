package portwatch

import (
	"context"
	"sync"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// WatchEvent represents a change observed during a watch cycle.
type WatchEvent struct {
	Port      scanner.Port
	State     string // "opened" or "closed"
	ObservedAt time.Time
}

// Watcher continuously polls open ports and emits WatchEvents on changes.
type Watcher struct {
	scanner  *scanner.Scanner
	interval time.Duration
	prev     map[string]scanner.Port
	mu       sync.Mutex
}

// New creates a Watcher with the given scanner and poll interval.
func New(s *scanner.Scanner, interval time.Duration) *Watcher {
	return &Watcher{
		scanner:  s,
		interval: interval,
		prev:     make(map[string]scanner.Port),
	}
}

// Watch starts the polling loop and sends events to the returned channel.
// The channel is closed when ctx is cancelled.
func (w *Watcher) Watch(ctx context.Context) <-chan WatchEvent {
	out := make(chan WatchEvent, 64)
	go func() {
		defer close(out)
		ticker := time.NewTicker(w.interval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				w.tick(ctx, out)
			}
		}
	}()
	return out
}

func (w *Watcher) tick(ctx context.Context, out chan<- WatchEvent) {
	ports, err := w.scanner.Scan(ctx)
	if err != nil {
		return
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	curr := make(map[string]scanner.Port, len(ports))
	for _, p := range ports {
		k := portKey(p)
		curr[k] = p
		if _, exists := w.prev[k]; !exists {
			select {
			case out <- WatchEvent{Port: p, State: "opened", ObservedAt: time.Now()}:
			default:
			}
		}
	}
	for k, p := range w.prev {
		if _, exists := curr[k]; !exists {
			select {
			case out <- WatchEvent{Port: p, State: "closed", ObservedAt: time.Now()}:
			default:
			}
		}
	}
	w.prev = curr
}

func portKey(p scanner.Port) string {
	return p.Proto + ":" + p.String()
}
