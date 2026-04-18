// Package portevict tracks ports that have been evicted (closed after being
// open for a very short duration), which may indicate ephemeral or flapping services.
package portevict

import (
	"sync"
	"time"

	"portwatch/internal/scanner"
)

// Entry records a single eviction event.
type Entry struct {
	Port     scanner.Port
	OpenedAt time.Time
	ClosedAt time.Time
	Duration time.Duration
}

// Tracker records ports closed within a minimum uptime threshold.
type Tracker struct {
	mu       sync.Mutex
	minUp    time.Duration
	evictions []Entry
}

// New creates a Tracker that flags ports closed before minUp has elapsed.
func New(minUp time.Duration) *Tracker {
	return &Tracker{minUp: minUp}
}

// Opened registers the time a port was first seen open.
func (t *Tracker) Opened(p scanner.Port, at time.Time) {
	// stored externally via Closed; we just need the pair
}

// Closed evaluates whether the port was evicted and records it if so.
func (t *Tracker) Closed(p scanner.Port, openedAt, closedAt time.Time) bool {
	d := closedAt.Sub(openedAt)
	if d < 0 {
		d = 0
	}
	if d >= t.minUp {
		return false
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	t.evictions = append(t.evictions, Entry{
		Port:     p,
		OpenedAt: openedAt,
		ClosedAt: closedAt,
		Duration: d,
	})
	return true
}

// All returns a copy of all recorded eviction entries.
func (t *Tracker) All() []Entry {
	t.mu.Lock()
	defer t.mu.Unlock()
	out := make([]Entry, len(t.evictions))
	copy(out, t.evictions)
	return out
}

// Reset clears all eviction history.
func (t *Tracker) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.evictions = nil
}
