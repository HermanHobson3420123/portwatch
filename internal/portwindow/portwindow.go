// Package portwindow tracks which ports were active during a sliding time window.
package portwindow

import (
	"sync"
	"time"

	"portwatch/internal/scanner"
)

// Entry holds the first and last seen timestamps for a port within the window.
type Entry struct {
	Port      scanner.Port
	FirstSeen time.Time
	LastSeen  time.Time
	Count     int
}

// Window tracks port activity over a configurable duration.
type Window struct {
	mu       sync.Mutex
	duration time.Duration
	entries  map[string]*Entry
	now      func() time.Time
}

// New creates a Window with the given duration.
func New(duration time.Duration) *Window {
	return &Window{
		duration: duration,
		entries:  make(map[string]*Entry),
		now:      time.Now,
	}
}

func portKey(p scanner.Port) string {
	return p.Proto + ":" + p.Addr
}

// Record marks a port as seen at the current time.
func (w *Window) Record(p scanner.Port) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.evict()
	key := portKey(p)
	if e, ok := w.entries[key]; ok {
		e.LastSeen = w.now()
		e.Count++
		return
	}
	now := w.now()
	w.entries[key] = &Entry{
		Port:      p,
		FirstSeen: now,
		LastSeen:  now,
		Count:     1,
	}
}

// Active returns all entries whose LastSeen is within the window duration.
func (w *Window) Active() []Entry {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.evict()
	out := make([]Entry, 0, len(w.entries))
	for _, e := range w.entries {
		out = append(out, *e)
	}
	return out
}

// evict removes entries that have not been seen within the window. Must be called with mu held.
func (w *Window) evict() {
	cutoff := w.now().Add(-w.duration)
	for k, e := range w.entries {
		if e.LastSeen.Before(cutoff) {
			delete(w.entries, k)
		}
	}
}

// Len returns the number of active entries.
func (w *Window) Len() int {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.evict()
	return len(w.entries)
}
