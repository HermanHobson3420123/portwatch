package suppress

import (
	"encoding/json"
	"errors"
	"os"
	"time"
)

// storedEntry is the JSON-serialisable form of a suppression entry.
type storedEntry struct {
	Port     uint16    `json:"port"`
	Protocol string    `json:"protocol"`
	Until    time.Time `json:"until"`
}

// SaveToFile persists all active suppression entries to a JSON file.
func (l *List) SaveToFile(path string) error {
	active := l.Active()
	stored := make([]storedEntry, len(active))
	for i, e := range active {
		stored[i] = storedEntry{Port: e.Port, Protocol: e.Protocol, Until: e.Until}
	}
	data, err := json.MarshalIndent(stored, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

// LoadFromFile reads suppression entries from a JSON file and adds them
// to the list. Entries that have already expired are silently skipped.
// Returns nil if the file does not exist.
func (l *List) LoadFromFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	var stored []storedEntry
	if err := json.Unmarshal(data, &stored); err != nil {
		return err
	}
	now := l.now()
	for _, e := range stored {
		if now.Before(e.Until) {
			l.Add(e.Port, e.Protocol, e.Until)
		}
	}
	return nil
}
