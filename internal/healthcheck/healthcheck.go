package healthcheck

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

// Status represents the current health of the daemon.
type Status struct {
	Healthy     bool      `json:"healthy"`
	LastScan    time.Time `json:"last_scan"`
	ScanCount   int64     `json:"scan_count"`
	AlertCount  int64     `json:"alert_count"`
	StartedAt   time.Time `json:"started_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Checker maintains and persists daemon health status.
type Checker struct {
	mu     sync.RWMutex
	status Status
	path   string
}

// New creates a Checker that writes status to path.
func New(path string) *Checker {
	return &Checker{
		path: path,
		status: Status{
			Healthy:   true,
			StartedAt: time.Now(),
		},
	}
}

// RecordScan updates the last scan timestamp and increments scan count.
func (c *Checker) RecordScan() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.status.LastScan = time.Now()
	c.status.ScanCount++
	c.status.UpdatedAt = time.Now()
	c.flush()
}

// RecordAlert increments the alert counter.
func (c *Checker) RecordAlert() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.status.AlertCount++
	c.status.UpdatedAt = time.Now()
	c.flush()
}

// SetHealthy marks the daemon healthy or unhealthy.
func (c *Checker) SetHealthy(ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.status.Healthy = ok
	c.status.UpdatedAt = time.Now()
	c.flush()
}

// Get returns a snapshot of the current status.
func (c *Checker) Get() Status {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.status
}

// flush writes status to disk; caller must hold mu.
func (c *Checker) flush() {
	data, err := json.MarshalIndent(c.status, "", "  ")
	if err != nil {
		return
	}
	_ = os.WriteFile(c.path, data, 0644)
}
