package snapshot

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Manager handles snapshot persistence with configurable directory and retention.
type Manager struct {
	dir       string
	retainDays int
}

// NewManager creates a Manager that stores snapshots in dir.
// retainDays specifies how many days of snapshots to keep (0 = keep all).
func NewManager(dir string, retainDays int) (*Manager, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("snapshot: create dir %s: %w", dir, err)
	}
	return &Manager{dir: dir, retainDays: retainDays}, nil
}

// LatestPath returns the path for the canonical latest snapshot file.
func (m *Manager) LatestPath() string {
	return filepath.Join(m.dir, "latest.json")
}

// ArchivePath returns a timestamped archive path for the current moment.
func (m *Manager) ArchivePath() string {
	stamp := time.Now().UTC().Format("20060102T150405Z")
	return filepath.Join(m.dir, fmt.Sprintf("snapshot_%s.json", stamp))
}

// Prune removes snapshot archive files older than retainDays.
// If retainDays is 0, no files are removed.
func (m *Manager) Prune() error {
	if m.retainDays == 0 {
		return nil
	}
	cutoff := time.Now().UTC().AddDate(0, 0, -m.retainDays)
	entries, err := os.ReadDir(m.dir)
	if err != nil {
		return fmt.Errorf("snapshot: read dir: %w", err)
	}
	for _, e := range entries {
		if e.Name() == "latest.json" {
			continue
		}
		info, err := e.Info()
		if err != nil {
			continue
		}
		if info.ModTime().Before(cutoff) {
			_ = os.Remove(filepath.Join(m.dir, e.Name()))
		}
	}
	return nil
}
