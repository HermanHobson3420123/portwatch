// Package suppress provides a mechanism to temporarily suppress
// alerts for specific ports, useful for planned maintenance windows.
package suppress

import (
	"sync"
	"time"
)

// Entry represents a single suppression rule.
type Entry struct {
	Port     uint16
	Protocol string
	Until    time.Time
}

// key uniquely identifies a port+protocol pair.
type key struct {
	port     uint16
	protocol string
}

// List holds active suppression entries.
type List struct {
	mu      sync.RWMutex
	entries map[key]time.Time
	now     func() time.Time
}

// New creates a new empty suppression list.
func New() *List {
	return &List{
		entries: make(map[key]time.Time),
		now:     time.Now,
	}
}

// Add suppresses alerts for the given port/protocol until the given time.
func (l *List) Add(port uint16, protocol string, until time.Time) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.entries[key{port, protocol}] = until
}

// Remove removes a suppression entry immediately.
func (l *List) Remove(port uint16, protocol string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.entries, key{port, protocol})
}

// IsSuppressed reports whether alerts for the given port/protocol
// should currently be suppressed.
func (l *List) IsSuppressed(port uint16, protocol string) bool {
	l.mu.RLock()
	defer l.mu.RUnlock()
	until, ok := l.entries[key{port, protocol}]
	if !ok {
		return false
	}
	return l.now().Before(until)
}

// Purge removes all expired suppression entries.
func (l *List) Purge() {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := l.now()
	for k, until := range l.entries {
		if !now.Before(until) {
			delete(l.entries, k)
		}
	}
}

// Active returns all currently active (non-expired) entries.
func (l *List) Active() []Entry {
	l.mu.RLock()
	defer l.mu.RUnlock()
	now := l.now()
	var out []Entry
	for k, until := range l.entries {
		if now.Before(until) {
			out = append(out, Entry{Port: k.port, Protocol: k.protocol, Until: until})
		}
	}
	return out
}
