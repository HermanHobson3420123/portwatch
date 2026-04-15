package suppress

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "suppress.json")

	now := time.Now()
	l := New()
	l.now = fixedNow(now)
	l.Add(80, "tcp", now.Add(1*time.Hour))
	l.Add(443, "tcp", now.Add(2*time.Hour))

	if err := l.SaveToFile(path); err != nil {
		t.Fatalf("SaveToFile: %v", err)
	}

	l2 := New()
	l2.now = fixedNow(now)
	if err := l2.LoadFromFile(path); err != nil {
		t.Fatalf("LoadFromFile: %v", err)
	}

	if !l2.IsSuppressed(80, "tcp") {
		t.Error("expected port 80/tcp to be suppressed after load")
	}
	if !l2.IsSuppressed(443, "tcp") {
		t.Error("expected port 443/tcp to be suppressed after load")
	}
}

func TestLoadMissingFileReturnsNil(t *testing.T) {
	l := New()
	err := l.LoadFromFile("/nonexistent/path/suppress.json")
	if err != nil {
		t.Errorf("expected nil for missing file, got %v", err)
	}
}

func TestLoadInvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	_ = os.WriteFile(path, []byte("not json"), 0o644)

	l := New()
	if err := l.LoadFromFile(path); err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestLoadSkipsExpiredEntries(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "suppress.json")
	now := time.Now()

	expired := []storedEntry{
		{Port: 8080, Protocol: "tcp", Until: now.Add(-1 * time.Minute)},
	}
	data, _ := json.MarshalIndent(expired, "", "  ")
	_ = os.WriteFile(path, data, 0o644)

	l := New()
	l.now = fixedNow(now)
	if err := l.LoadFromFile(path); err != nil {
		t.Fatalf("LoadFromFile: %v", err)
	}
	if l.IsSuppressed(8080, "tcp") {
		t.Error("expected expired entry to be skipped on load")
	}
}
