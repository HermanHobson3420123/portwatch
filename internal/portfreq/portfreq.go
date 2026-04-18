package portfreq

import (
	"sync"
	"time"

	"portwatch/internal/scanner"
)

// Entry records how often a port has been seen open.
type Entry struct {
	Port      scanner.Port
	SeenCount int
	FirstSeen time.Time
	LastSeen  time.Time
}

// Tracker counts scan appearances per port.
type Tracker struct {
	mu      sync.Mutex
	entries map[string]*Entry
}

func portKey(p scanner.Port) string {
	return p.Proto + ":" + p.Addr
}

// New returns an empty Tracker.
func New() *Tracker {
	return &Tracker{entries: make(map[string]*Entry)}
}

// Record marks the port as seen at t.
func (tr *Tracker) Record(p scanner.Port, t time.Time) {
	tr.mu.Lock()
	defer tr.mu.Unlock()
	k := portKey(p)
	e, ok := tr.entries[k]
	if !ok {
		tr.entries[k] = &Entry{Port: p, SeenCount: 1, FirstSeen: t, LastSeen: t}
		return
	}
	e.SeenCount++
	e.LastSeen = t
}

// Top returns the n most-frequently seen ports.
func (tr *Tracker) Top(n int) []Entry {
	tr.mu.Lock()
	defer tr.mu.Unlock()
	all := make([]Entry, 0, len(tr.entries))
	for _, e := range tr.entries {
		all = append(all, *e)
	}
	sortByCount(all)
	if n > 0 && n < len(all) {
		return all[:n]
	}
	return all
}

// Reset clears all recorded data.
func (tr *Tracker) Reset() {
	tr.mu.Lock()
	defer tr.mu.Unlock()
	tr.entries = make(map[string]*Entry)
}

func sortByCount(es []Entry) {
	for i := 1; i < len(es); i++ {
		for j := i; j > 0 && es[j].SeenCount > es[j-1].SeenCount; j-- {
			es[j], es[j-1] = es[j-1], es[j]
		}
	}
}
