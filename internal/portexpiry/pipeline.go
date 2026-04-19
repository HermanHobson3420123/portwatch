package portexpiry

import (
	"context"
	"time"

	"portwatch/internal/scanner"
)

// ExpiryEvent represents a port that has expired (not seen within TTL).
type ExpiryEvent struct {
	Port      scanner.Port
	LastSeen  time.Time
	ExpiredAt time.Time
}

// NewPipeline starts a goroutine that periodically evicts expired ports and
// sends ExpiryEvents on the returned channel. It also feeds observed ports
// into the tracker via the input channel.
func NewPipeline(ctx context.Context, in <-chan scanner.Port, ttl, interval time.Duration) <-chan ExpiryEvent {
	out := make(chan ExpiryEvent, 16)
	tr := New(ttl)
	go func() {
		defer close(out)
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case port, ok := <-in:
				if !ok {
					return
				}
				tr.Seen(port, time.Now())
			case now := <-ticker.C:
				for _, e := range tr.Evict(now) {
					select {
					case out <- ExpiryEvent{Port: e.Port, LastSeen: e.LastSeen, ExpiredAt: now}:
					case <-ctx.Done():
						return
					}
				}
			}
		}
	}()
	return out
}
