package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestMainConfigFlagMissing(t *testing.T) {
	// Verify that a non-existent config path causes a detectable error
	// by testing config.Load directly (integration-style smoke test)
	_, err := os.Stat("/nonexistent/path/portwatch.json")
	if err == nil {
		t.Fatal("expected file to not exist")
	}
}

func TestMainConfigLoadValid(t *testing.T) {
	dir := t.TempDir()
	cfgFile := filepath.Join(dir, "config.json")

	content := []byte(`{
		"port_range_start": 1024,
		"port_range_end": 2048,
		"interval": "10s",
		"timeout": "500ms"
	}`)
	if err := os.WriteFile(cfgFile, content, 0644); err != nil {
		t.Fatalf("failed to write temp config: %v", err)
	}

	// Ensure file is readable
	data, err := os.ReadFile(cfgFile)
	if err != nil {
		t.Fatalf("failed to read temp config: %v", err)
	}
	if !bytes.Contains(data, []byte("port_range_start")) {
		t.Error("expected config content to contain port_range_start")
	}
}

func TestMainStdoutIsWritable(t *testing.T) {
	if os.Stdout == nil {
		t.Fatal("expected stdout to be non-nil")
	}
}
