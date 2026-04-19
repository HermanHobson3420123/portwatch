package portexpiry

import (
	"sync"
	"time"

	"portwatch/internal/scanner"
)

// Entry tracks when a port was last seen open.
type Entry struct {
	Port     scanner.Port
	LastSeen time.Time
	TTL      time.Duration
}

// Expired returns true if the entry has exceeded its TTL.
func (e Entry) Expired(now time.Time) bool {
	return now.Sub(e.LastSeen) > e.TTL
}

// Tracker monitors ports and emits those that have not been seen within TTL.
type Tracker struct {
	mu      sync.Mutex
	entries map[string]Entry
	ttl     time.Duration
}

func portKey(p scanner.Port) string {
	return p.Proto + ":" + p.String()
}

// New creates a Tracker with the given default TTL.
func New(ttl time.Duration) *Tracker {
	return &Tracker{
		entries: make(map[string]Entry),
		ttl:     ttl,
	}
}

// Seen records that a port was observed at the given time.
func (t *Tracker) Seen(p scanner.Port, now time.Time) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.entries[portKey(p)] = Entry{Port: p, LastSeen: now, TTL: t.ttl}
}

// Expired returns all entries whose TTL has elapsed relative to now.
func (t *Tracker) Expired(now time.Time) []Entry {
	t.mu.Lock()
	defer t.mu.Unlock()
	var out []Entry
	for _, e := range t.entries {
		if e.Expired(now) {
			out = append(out, e)
		}
	}
	return out
}

// Evict removes expired entries and returns them.
func (t *Tracker) Evict(now time.Time) []Entry {
	t.mu.Lock()
	defer t.mu.Unlock()
	var out []Entry
	for k, e := range t.entries {
		if e.Expired(now) {
			out = append(out, e)
			delete(t.entries, k)
		}
	}
	return out
}

// Len returns the number of tracked entries.
func (t *Tracker) Len() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return len(t.entries)
}
