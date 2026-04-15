package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Config holds the full portwatch configuration.
type Config struct {
	Ports        []int         `json:"ports"`
	Protocols    []string      `json:"protocols"`
	Interval     time.Duration `json:"interval"`
	AlertLog     string        `json:"alert_log"`
	SnapshotDir  string        `json:"snapshot_dir"`
	RetainDays   int           `json:"retain_days"`
}

// DefaultConfig returns a Config populated with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Ports:       []int{},
		Protocols:   []string{"tcp"},
		Interval:    30 * time.Second,
		AlertLog:    "",
		SnapshotDir: ".portwatch/snapshots",
		RetainDays:  30,
	}
}

// rawConfig mirrors Config but uses a plain int for JSON duration parsing.
type rawConfig struct {
	Ports       []int    `json:"ports"`
	Protocols   []string `json:"protocols"`
	IntervalSec int      `json:"interval_seconds"`
	AlertLog    string   `json:"alert_log"`
	SnapshotDir string   `json:"snapshot_dir"`
	RetainDays  int      `json:"retain_days"`
}

// Load reads a JSON config file and merges it over the defaults.
func Load(path string) (Config, error) {
	cfg := DefaultConfig()
	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, fmt.Errorf("config: read %s: %w", path, err)
	}
	var raw rawConfig
	if err := json.Unmarshal(data, &raw); err != nil {
		return cfg, fmt.Errorf("config: parse %s: %w", path, err)
	}
	if len(raw.Ports) > 0 {
		cfg.Ports = raw.Ports
	}
	if len(raw.Protocols) > 0 {
		cfg.Protocols = raw.Protocols
	}
	if raw.IntervalSec > 0 {
		cfg.Interval = time.Duration(raw.IntervalSec) * time.Second
	}
	if raw.AlertLog != "" {
		cfg.AlertLog = raw.AlertLog
	}
	if raw.SnapshotDir != "" {
		cfg.SnapshotDir = raw.SnapshotDir
	}
	if raw.RetainDays > 0 {
		cfg.RetainDays = raw.RetainDays
	}
	return cfg, nil
}
