package portmap

import (
	"fmt"
	"sync"

	"github.com/user/portwatch/internal/scanner"
)

// Map maintains a thread-safe map of port key -> Port for quick lookup.
type Map struct {
	mu    sync.RWMutex
	ports map[string]scanner.Port
}

// New returns an empty Map.
func New() *Map {
	return &Map{ports: make(map[string]scanner.Port)}
}

// Set inserts or updates a port entry.
func (m *Map) Set(p scanner.Port) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.ports[key(p)] = p
}

// Delete removes a port entry if present.
func (m *Map) Delete(p scanner.Port) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.ports, key(p))
}

// Has returns true when the port is currently tracked.
func (m *Map) Has(p scanner.Port) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	_, ok := m.ports[key(p)]
	return ok
}

// All returns a snapshot slice of all tracked ports.
func (m *Map) All() []scanner.Port {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]scanner.Port, 0, len(m.ports))
	for _, p := range m.ports {
		out = append(out, p)
	}
	return out
}

// Len returns the number of tracked ports.
func (m *Map) Len() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.ports)
}

func key(p scanner.Port) string {
	return fmt.Sprintf("%s:%d", p.Proto, p.Number)
}
