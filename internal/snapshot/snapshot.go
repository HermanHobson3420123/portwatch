package snapshot

import (
	"encoding/json"
	"os"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Snapshot holds a recorded set of open ports at a point in time.
type Snapshot struct {
	Timestamp time.Time      `json:"timestamp"`
	Ports     []scanner.Port `json:"ports"`
}

// Save writes the snapshot to the given file path as JSON.
func Save(path string, ports []scanner.Port) error {
	snap := Snapshot{
		Timestamp: time.Now().UTC(),
		Ports:     ports,
	}
	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// Load reads a snapshot from the given file path.
func Load(path string) (*Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return nil, err
	}
	return &snap, nil
}

// Diff compares two port slices and returns opened and closed ports.
func Diff(previous, current []scanner.Port) (opened, closed []scanner.Port) {
	prevSet := make(map[string]struct{}, len(previous))
	for _, p := range previous {
		prevSet[p.String()] = struct{}{}
	}
	currSet := make(map[string]struct{}, len(current))
	for _, p := range current {
		currSet[p.String()] = struct{}{}
	}
	for _, p := range current {
		if _, found := prevSet[p.String()]; !found {
			opened = append(opened, p)
		}
	}
	for _, p := range previous {
		if _, found := currSet[p.String()]; !found {
			closed = append(closed, p)
		}
	}
	return opened, closed
}
