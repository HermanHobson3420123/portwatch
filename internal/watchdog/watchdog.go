package watchdog

import (
	"context"
	"log"
	"time"

	"portwatch/internal/scanner"
)

// HealthStatus represents the current health of the watchdog.
type HealthStatus struct {
	LastScan  time.Time
	ScanCount int64
	Healthy   bool
}

// Watchdog monitors the scanner loop and restarts it if it stalls.
type Watchdog struct {
	interval  time.Duration
	timeout   time.Duration
	scanner   *scanner.Scanner
	heartbeat chan struct{}
	status    HealthStatus
}

// New creates a Watchdog that expects a heartbeat every timeout duration.
func New(s *scanner.Scanner, interval, timeout time.Duration) *Watchdog {
	return &Watchdog{
		interval:  interval,
		timeout:   timeout,
		scanner:   s,
		heartbeat: make(chan struct{}, 1),
	}
}

// Beat records that the scanner completed a scan cycle.
func (w *Watchdog) Beat() {
	w.status.LastScan = time.Now()
	w.status.ScanCount++
	select {
	case w.heartbeat <- struct{}{}:
	default:
	}
}

// Status returns a snapshot of the current health status.
func (w *Watchdog) Status() HealthStatus {
	return w.status
}

// Watch runs the watchdog loop, logging warnings when heartbeats are missed.
func (w *Watchdog) Watch(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if w.status.LastScan.IsZero() {
				continue
			}
			staleDuration := time.Since(w.status.LastScan)
			if staleDuration > w.timeout {
				w.status.Healthy = false
				log.Printf("[watchdog] WARNING: no scan heartbeat for %s (last: %s)",
					staleDuration.Round(time.Second),
					w.status.LastScan.Format(time.RFC3339),
				)
			} else {
				w.status.Healthy = true
			}
		case <-w.heartbeat:
			w.status.Healthy = true
		}
	}
}
