package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/snapshot"
)

func makePorts(nums ...int) []scanner.Port {
	ports := make([]scanner.Port, len(nums))
	for i, n := range nums {
		ports[i] = scanner.Port{Number: n, Protocol: "tcp"}
	}
	return ports
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")
	ports := makePorts(80, 443, 8080)

	if err := snapshot.Save(path, ports); err != nil {
		t.Fatalf("Save failed: %v", err)
	}
	snap, err := snapshot.Load(path)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if len(snap.Ports) != len(ports) {
		t.Errorf("expected %d ports, got %d", len(ports), len(snap.Ports))
	}
	if snap.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestLoadMissingFile(t *testing.T) {
	_, err := snapshot.Load("/nonexistent/path/snap.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestLoadInvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	_ = os.WriteFile(path, []byte("not json{"), 0644)
	_, err := snapshot.Load(path)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestDiffOpened(t *testing.T) {
	prev := makePorts(80, 443)
	curr := makePorts(80, 443, 8080)
	opened, closed := snapshot.Diff(prev, curr)
	if len(opened) != 1 || opened[0].Number != 8080 {
		t.Errorf("expected port 8080 opened, got %v", opened)
	}
	if len(closed) != 0 {
		t.Errorf("expected no closed ports, got %v", closed)
	}
}

func TestDiffClosed(t *testing.T) {
	prev := makePorts(80, 443, 8080)
	curr := makePorts(80, 443)
	opened, closed := snapshot.Diff(prev, curr)
	if len(closed) != 1 || closed[0].Number != 8080 {
		t.Errorf("expected port 8080 closed, got %v", closed)
	}
	if len(opened) != 0 {
		t.Errorf("expected no opened ports, got %v", opened)
	}
}

func TestDiffNoChange(t *testing.T) {
	ports := makePorts(80, 443)
	opened, closed := snapshot.Diff(ports, ports)
	if len(opened) != 0 || len(closed) != 0 {
		t.Errorf("expected no diff, got opened=%v closed=%v", opened, closed)
	}
}

func TestDiffBothOpenedAndClosed(t *testing.T) {
	prev := makePorts(80, 443, 8080)
	curr := makePorts(80, 443, 9090)
	opened, closed := snapshot.Diff(prev, curr)
	if len(opened) != 1 || opened[0].Number != 9090 {
		t.Errorf("expected port 9090 opened, got %v", opened)
	}
	if len(closed) != 1 || closed[0].Number != 8080 {
		t.Errorf("expected port 8080 closed, got %v", closed)
	}
}
