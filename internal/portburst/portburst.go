// Package portburst detects rapid sequences of port open/close events
// ("bursts") within a sliding time window and emits a summary when the
// burst threshold is exceeded.
package portburst

import (
	"fmt"
	"sync"
	"time"

	"portwatch/internal/scanner"
)

// Event is emitted when a burst is detected.
type Event struct {
	Port      scanner.Port
	Count     int
	Window    time.Duration
	DetectedAt time.Time
}

func (e Event) String() string {
	return fmt.Sprintf("burst detected on %s: %d events in %s",
		e.Port, e.Count, e.Window)
}

// Detector tracks port activity timestamps and fires when the count of
// events within Window exceeds Threshold.
type Detector struct {
	mu        sync.Mutex
	Window    time.Duration
	Threshold int
	clock     func() time.Time
	events    map[string][]time.Time // portKey -> timestamps
}

func portKey(p scanner.Port) string {
	return fmt.Sprintf("%s:%d", p.Proto, p.Number)
}

// New returns a Detector with the given window and threshold.
// A zero or negative threshold defaults to 5.
func New(window time.Duration, threshold int) *Detector {
	if threshold <= 0 {
		threshold = 5
	}
	return &Detector{
		Window:    window,
		Threshold: threshold,
		clock:     time.Now,
		events:    make(map[string][]time.Time),
	}
}

// Record registers an activity event for the given port and returns a
// non-nil *Event if the burst threshold has been crossed.
func (d *Detector) Record(p scanner.Port) *Event {
	d.mu.Lock()
	defer d.mu.Unlock()

	now := d.clock()
	cutoff := now.Add(-d.Window)
	key := portKey(p)

	ts := d.events[key]
	// prune old entries
	filtered := ts[:0]
	for _, t := range ts {
		if t.After(cutoff) {
			filtered = append(filtered, t)
		}
	}
	filtered = append(filtered, now)
	d.events[key] = filtered

	if len(filtered) >= d.Threshold {
		return &Event{
			Port:       p,
			Count:      len(filtered),
			Window:     d.Window,
			DetectedAt: now,
		}
	}
	return nil
}

// Reset clears all recorded events for a port.
func (d *Detector) Reset(p scanner.Port) {
	d.mu.Lock()
	defer d.mu.Unlock()
	delete(d.events, portKey(p))
}
