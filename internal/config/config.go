package config

import (
	"encoding/json"
	"os"
	"time"
)

// Config holds the runtime configuration for portwatch.
type Config struct {
	// Ports to monitor. If empty, all ports are monitored.
	Ports []uint16 `json:"ports"`

	// Interval between scans.
	Interval time.Duration `json:"interval"`

	// Host to scan (defaults to localhost).
	Host string `json:"host"`

	// AlertCommand is an optional shell command to run on changes.
	AlertCommand string `json:"alert_command"`
}

// DefaultConfig returns a Config populated with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		Ports:    []uint16{},
		Interval: 5 * time.Second,
		Host:     "127.0.0.1",
	}
}

// Load reads a JSON config file from path and returns a Config.
// Missing fields fall back to defaults.
func Load(path string) (*Config, error) {
	cfg := DefaultConfig()

	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, err
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	if err := dec.Decode(cfg); err != nil {
		return nil, err
	}

	if cfg.Host == "" {
		cfg.Host = "127.0.0.1"
	}
	if cfg.Interval <= 0 {
		cfg.Interval = 5 * time.Second
	}

	return cfg, nil
}
