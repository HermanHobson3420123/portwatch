package config

import (
	"encoding/json"
	"os"
	"testing"
	"time"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.Host != "127.0.0.1" {
		t.Errorf("expected default host 127.0.0.1, got %s", cfg.Host)
	}
	if cfg.Interval != 5*time.Second {
		t.Errorf("expected default interval 5s, got %v", cfg.Interval)
	}
	if len(cfg.Ports) != 0 {
		t.Errorf("expected empty default ports, got %v", cfg.Ports)
	}
}

func TestLoadMissingFile(t *testing.T) {
	cfg, err := Load("/tmp/portwatch_nonexistent_config.json")
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if cfg == nil {
		t.Fatal("expected non-nil config")
	}
	if cfg.Host != "127.0.0.1" {
		t.Errorf("expected default host, got %s", cfg.Host)
	}
}

func TestLoadValidConfig(t *testing.T) {
	data := map[string]interface{}{
		"ports":    []int{80, 443, 8080},
		"interval": "10s",
		"host":     "0.0.0.0",
	}

	f, err := os.CreateTemp("", "portwatch-config-*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())

	if err := json.NewEncoder(f).Encode(data); err != nil {
		t.Fatal(err)
	}
	f.Close()

	cfg, err := Load(f.Name())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Host != "0.0.0.0" {
		t.Errorf("expected host 0.0.0.0, got %s", cfg.Host)
	}
	if len(cfg.Ports) != 3 {
		t.Errorf("expected 3 ports, got %d", len(cfg.Ports))
	}
	if cfg.Interval != 10*time.Second {
		t.Errorf("expected 10s interval, got %v", cfg.Interval)
	}
}

func TestLoadInvalidJSON(t *testing.T) {
	f, err := os.CreateTemp("", "portwatch-bad-*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	f.WriteString("not valid json")
	f.Close()

	_, err = Load(f.Name())
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}
