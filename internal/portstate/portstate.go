package portstate

import (
	"sync"
	"time"

	"portwatch/internal/scanner"
)

// State represents the current known state of a single port.
type State struct {
	Port      scanner.Port
	OpenSince time.Time
	LastSeen  time.Time
	OpenCount int
}

// Tracker maintains a live map of currently open ports.
type Tracker struct {
	mu     sync.RWMutex
	states map[string]*State
}

func portKey(p scanner.Port) string {
	return p.Proto + ":" + p.String()
}

// New returns an initialised Tracker.
func New() *Tracker {
	return &Tracker{states: make(map[string]*State)}
}

// Open records a port as open. If already tracked it updates LastSeen.
func (t *Tracker) Open(p scanner.Port, now time.Time) {
	t.mu.Lock()
	defer t.mu.Unlock()
	key := portKey(p)
	if s, ok := t.states[key]; ok {
		s.LastSeen = now
		return
	}
	t.states[key] = &State{Port: p, OpenSince: now, LastSeen: now, OpenCount: 1}
}

// Close removes a port from the tracker.
func (t *Tracker) Close(p scanner.Port) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.states, portKey(p))
}

// Get returns the State for a port and whether it exists.
func (t *Tracker) Get(p scanner.Port) (State, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	s, ok := t.states[portKey(p)]
	if !ok {
		return State{}, false
	}
	return *s, true
}

// All returns a snapshot of all tracked states.
func (t *Tracker) All() []State {
	t.mu.RLock()
	defer t.mu.RUnlock()
	out := make([]State, 0, len(t.states))
	for _, s := range t.states {
		out = append(out, *s)
	}
	return out
}

// Len returns the number of currently tracked open ports.
func (t *Tracker) Len() int {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return len(t.states)
}
