package portmap

import (
	"github.com/user/portwatch/internal/scanner"
)

// Delta describes ports that were opened or closed between two syncs.
type Delta struct {
	Opened []scanner.Port
	Closed []scanner.Port
}

// Sync updates the Map to match current and returns the Delta.
// Ports in current but not in the map are Opened.
// Ports in the map but not in current are Closed.
func (m *Map) Sync(current []scanner.Port) Delta {
	currentSet := make(map[string]scanner.Port, len(current))
	for _, p := range current {
		currentSet[key(p)] = p
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	var delta Delta

	for k, p := range m.ports {
		if _, ok := currentSet[k]; !ok {
			delta.Closed = append(delta.Closed, p)
			delete(m.ports, k)
		}
	}

	for k, p := range currentSet {
		if _, ok := m.ports[k]; !ok {
			delta.Opened = append(delta.Opened, p)
			m.ports[k] = p
		}
	}

	return delta
}
