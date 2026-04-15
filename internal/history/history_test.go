package history_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"portwatch/internal/history"
	"portwatch/internal/scanner"
)

func makePort(proto string, num int) scanner.Port {
	return scanner.Port{Protocol: proto, Number: num}
}

func TestAppendAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "history.json")

	h := &history.History{}
	h.Append(history.EventOpened, makePort("tcp", 8080))
	h.Append(history.EventClosed, makePort("udp", 53))

	if len(h.Events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(h.Events))
	}

	if err := history.Save(path, h); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := history.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(loaded.Events) != 2 {
		t.Fatalf("expected 2 loaded events, got %d", len(loaded.Events))
	}
	if loaded.Events[0].Type != history.EventOpened {
		t.Errorf("expected opened, got %s", loaded.Events[0].Type)
	}
	if loaded.Events[1].Port.Number != 53 {
		t.Errorf("expected port 53, got %d", loaded.Events[1].Port.Number)
	}
}

func TestLoadMissingFileReturnsEmpty(t *testing.T) {
	h, err := history.Load("/nonexistent/path/history.json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(h.Events) != 0 {
		t.Errorf("expected empty history, got %d events", len(h.Events))
	}
}

func TestLoadInvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	_ = os.WriteFile(path, []byte("not json{"), 0o644)

	_, err := history.Load(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

func TestPruneRemovesOldEvents(t *testing.T) {
	h := &history.History{}
	old := history.Event{
		Timestamp: time.Now().UTC().Add(-2 * time.Hour),
		Type:      history.EventOpened,
		Port:      makePort("tcp", 22),
	}
	recent := history.Event{
		Timestamp: time.Now().UTC().Add(-1 * time.Minute),
		Type:      history.EventClosed,
		Port:      makePort("tcp", 443),
	}
	h.Events = []history.Event{old, recent}

	trimmed := history.Prune(h, 30*time.Minute)
	if len(trimmed.Events) != 1 {
		t.Fatalf("expected 1 event after prune, got %d", len(trimmed.Events))
	}
	if trimmed.Events[0].Port.Number != 443 {
		t.Errorf("expected port 443, got %d", trimmed.Events[0].Port.Number)
	}
}

func TestSaveProducesValidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "history.json")

	h := &history.History{}
	h.Append(history.EventOpened, makePort("tcp", 9090))

	if err := history.Save(path, h); err != nil {
		t.Fatalf("Save: %v", err)
	}

	data, _ := os.ReadFile(path)
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
}
