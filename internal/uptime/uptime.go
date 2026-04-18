package uptime

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

// Record tracks how long a port has been continuously open.
type Record struct {
	Protocol string    `json:"protocol"`
	Port     uint16    `json:"port"`
	FirstSeen time.Time `json:"first_seen"`
	LastSeen  time.Time `json:"last_seen"`
}

// Duration returns how long the port has been open.
func (r Record) Duration() time.Duration {
	return r.LastSeen.Sub(r.FirstSeen)
}

// Tracker maintains uptime records for open ports.
type Tracker struct {
	mu      sync.Mutex
	records map[string]*Record
	path    string
}

func key(protocol string, port uint16) string {
	return protocol + ":" + string(rune(port))
}

// New creates a Tracker, loading persisted state from path if present.
func New(path string) (*Tracker, error) {
	t := &Tracker{path: path, records: make(map[string]*Record)}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return t, nil
	}
	if err != nil {
		return nil, err
	}
	var records []*Record
	if err := json.Unmarshal(data, &records); err != nil {
		return nil, err
	}
	for _, r := range records {
		t.records[key(r.Protocol, r.Port)] = r
	}
	return t, nil
}

// Opened registers a port as open at the given time.
func (t *Tracker) Opened(protocol string, port uint16, at time.Time) {
	t.mu.Lock()
	defer t.mu.Unlock()
	k := key(protocol, port)
	if _, ok := t.records[k]; !ok {
		t.records[k] = &Record{Protocol: protocol, Port: port, FirstSeen: at, LastSeen: at}
	}
}

// Seen updates the last-seen timestamp for an open port.
func (t *Tracker) Seen(protocol string, port uint16, at time.Time) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if r, ok := t.records[key(protocol, port)]; ok {
		r.LastSeen = at
	}
}

// Closed removes a port from the tracker.
func (t *Tracker) Closed(protocol string, port uint16) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.records, key(protocol, port))
}

// Get returns the uptime record for a port, if present.
func (t *Tracker) Get(protocol string, port uint16) (*Record, bool) {
	t.mu.Lock()
	defer t.mu.Unlock()
	r, ok := t.records[key(protocol, port)]
	return r, ok
}

// Save persists the current records to disk.
func (t *Tracker) Save() error {
	t.mu.Lock()
	defer t.mu.Unlock()
	records := make([]*Record, 0, len(t.records))
	for _, r := range t.records {
		records = append(records, r)
	}
	data, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(t.path, data, 0o644)
}
