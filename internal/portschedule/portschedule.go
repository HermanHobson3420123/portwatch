package portschedule

import (
	"sync"
	"time"
)

// Entry holds a scheduled scan window for a port range.
type Entry struct {
	Label    string
	Start    time.Time
	End      time.Time
	PortLow  int
	PortHigh int
	Protocol string
}

// Active returns true if the entry is currently within its scheduled window.
func (e Entry) Active(now time.Time) bool {
	return !now.Before(e.Start) && now.Before(e.End)
}

// Schedule manages a collection of scan window entries.
type Schedule struct {
	mu      sync.RWMutex
	entries []Entry
}

// New returns an empty Schedule.
func New() *Schedule {
	return &Schedule{}
}

// Add inserts a new scheduled window entry.
func (s *Schedule) Add(e Entry) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries = append(s.entries, e)
}

// Remove deletes all entries matching the given label.
func (s *Schedule) Remove(label string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := s.entries[:0]
	for _, e := range s.entries {
		if e.Label != label {
			out = append(out, e)
		}
	}
	s.entries = out
}

// ActiveNow returns all entries active at the given time.
func (s *Schedule) ActiveNow(now time.Time) []Entry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var out []Entry
	for _, e := range s.entries {
		if e.Active(now) {
			out = append(out, e)
		}
	}
	return out
}

// All returns a copy of all entries.
func (s *Schedule) All() []Entry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Entry, len(s.entries))
	copy(out, s.entries)
	return out
}
