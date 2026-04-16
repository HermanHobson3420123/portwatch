package metrics

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

// Snapshot holds a point-in-time view of runtime counters.
type Snapshot struct {
	ScansTotal  int64     `json:"scans_total"`
	AlertsTotal int64     `json:"alerts_total"`
	OpenedTotal int64     `json:"opened_total"`
	ClosedTotal int64     `json:"closed_total"`
	LastScan    time.Time `json:"last_scan"`
	LastAlert   time.Time `json:"last_alert"`
}

// Collector accumulates runtime counters in a thread-safe manner.
type Collector struct {
	mu sync.Mutex
	snap Snapshot
}

// New returns a zeroed Collector.
func New() *Collector { return &Collector{} }

func (c *Collector) RecordScan() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.snap.ScansTotal++
	c.snap.LastScan = time.Now()
}

func (c *Collector) RecordOpened() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.snap.AlertsTotal++
	c.snap.OpenedTotal++
	c.snap.LastAlert = time.Now()
}

func (c *Collector) RecordClosed() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.snap.AlertsTotal++
	c.snap.ClosedTotal++
	c.snap.LastAlert = time.Now()
}

// Snapshot returns a copy of the current counters.
func (c *Collector) Snapshot() Snapshot {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.snap
}

// Save writes the current snapshot as JSON to path.
func (c *Collector) Save(path string) error {
	c.mu.Lock()
	snap := c.snap
	c.mu.Unlock()
	b, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0o644)
}

// Load reads a previously saved snapshot from path.
func Load(path string) (Snapshot, error) {
	var s Snapshot
	b, err := os.ReadFile(path)
	if err != nil {
		return s, err
	}
	err = json.Unmarshal(b, &s)
	return s, err
}
