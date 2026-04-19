package portcount

import (
	"fmt"
	"sync"

	"portwatch/internal/scanner"
)

// Snapshot holds a count summary for a set of ports.
type Snapshot struct {
	Total int `json:"total"`
	TCP   int `json:"tcp"`
	UDP   int `json:"udp"`
}

// Counter tracks the number of open ports over time.
type Counter struct {
	mu   sync.Mutex
	last Snapshot
}

// New returns a new Counter.
func New() *Counter {
	return &Counter{}
}

// Record updates the counter from the latest port scan result.
func (c *Counter) Record(ports []scanner.Port) Snapshot {
	c.mu.Lock()
	defer c.mu.Unlock()

	snap := Snapshot{Total: len(ports)}
	for _, p := range ports {
		switch p.Protocol {
		case "tcp":
			snap.TCP++
		case "udp":
			snap.UDP++
		}
	}
	c.last = snap
	return snap
}

// Last returns the most recently recorded snapshot.
func (c *Counter) Last() Snapshot {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.last
}

// Summary returns a human-readable summary string.
func (c *Counter) Summary() string {
	s := c.Last()
	return fmt.Sprintf("open ports: total=%d tcp=%d udp=%d", s.Total, s.TCP, s.UDP)
}
