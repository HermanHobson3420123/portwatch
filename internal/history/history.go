package history

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"portwatch/internal/scanner"
)

// EventType describes whether a port was opened or closed.
type EventType string

const (
	EventOpened EventType = "opened"
	EventClosed EventType = "closed"
)

// Event records a single port change.
type Event struct {
	Timestamp time.Time    `json:"timestamp"`
	Type      EventType    `json:"type"`
	Port      scanner.Port `json:"port"`
}

// History holds an ordered list of port change events.
type History struct {
	Events []Event `json:"events"`
}

// Append adds a new event to the history.
func (h *History) Append(t EventType, p scanner.Port) {
	h.Events = append(h.Events, Event{
		Timestamp: time.Now().UTC(),
		Type:      t,
		Port:      p,
	})
}

// Save writes the history to the given file path as JSON.
func Save(path string, h *History) error {
	data, err := json.MarshalIndent(h, "", "  ")
	if err != nil {
		return fmt.Errorf("history: marshal: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("history: write %s: %w", path, err)
	}
	return nil
}

// Load reads history from the given file path.
// Returns an empty History if the file does not exist.
func Load(path string) (*History, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return &History{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("history: read %s: %w", path, err)
	}
	var h History
	if err := json.Unmarshal(data, &h); err != nil {
		return nil, fmt.Errorf("history: unmarshal: %w", err)
	}
	return &h, nil
}

// Prune removes events older than the given duration, returning the trimmed History.
func Prune(h *History, maxAge time.Duration) *History {
	cutoff := time.Now().UTC().Add(-maxAge)
	var kept []Event
	for _, e := range h.Events {
		if e.Timestamp.After(cutoff) {
			kept = append(kept, e)
		}
	}
	return &History{Events: kept}
}
