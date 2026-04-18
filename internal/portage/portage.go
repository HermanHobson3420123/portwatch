// Package portage tracks how long each port has been in its current state.
package portage

import (
	"fmt"
	"sync"
	"time"

	"portwatch/internal/scanner"
)

// Entry holds the first-seen time for a port.
type Entry struct {
	Port      scanner.Port
	FirstSeen time.Time
}

// Age returns how long the port has been open.
func (e Entry) Age(now time.Time) time.Duration {
	return now.Sub(e.FirstSeen)
}

// AgeString returns a human-readable age string.
func (e Entry) AgeString(now time.Time) string {
	d := e.Age(now)
	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	}
	if d < time.Hour {
		return fmt.Sprintf("%dm", int(d.Minutes()))
	}
	return fmt.Sprintf("%dh", int(d.Hours()))
}

func portKey(p scanner.Port) string {
	return fmt.Sprintf("%s:%d", p.Proto, p.Port)
}

// Tracker records when each port was first seen open.
type Tracker struct {
	mu      sync.Mutex
	entries map[string]Entry
	now     func() time.Time
}

// New returns a new Tracker.
func New() *Tracker {
	return &Tracker{
		entries: make(map[string]Entry),
		now:     time.Now,
	}
}

// Opened records a port as newly opened.
func (t *Tracker) Opened(p scanner.Port) {
	t.mu.Lock()
	defer t.mu.Unlock()
	key := portKey(p)
	if _, ok := t.entries[key]; !ok {
		t.entries[key] = Entry{Port: p, FirstSeen: t.now()}
	}
}

// Closed removes a port from tracking.
func (t *Tracker) Closed(p scanner.Port) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.entries, portKey(p))
}

// Get returns the entry for a port and whether it exists.
func (t *Tracker) Get(p scanner.Port) (Entry, bool) {
	t.mu.Lock()
	defer t.mu.Unlock()
	e, ok := t.entries[portKey(p)]
	return e, ok
}

// All returns all tracked entries sorted by first-seen ascending.
func (t *Tracker) All() []Entry {
	t.mu.Lock()
	defer t.mu.Unlock()
	out := make([]Entry, 0, len(t.entries))
	for _, e := range t.entries {
		out = append(out, e)
	}
	return out
}
