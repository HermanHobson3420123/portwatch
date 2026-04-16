package audit_test

import (
	"os"
	"path/filepath"
	"testing"

	"portwatch/internal/audit"
	"portwatch/internal/scanner"
)

func makePort(proto string, port uint16) scanner.Port {
	return scanner.Port{Proto: proto, Port: port, State: "open"}
}

func TestAppendAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "audit.json")

	l, err := audit.New(path)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := l.Append(audit.KindOpened, makePort("tcp", 80), ""); err != nil {
		t.Fatalf("Append: %v", err)
	}

	l2, err := audit.New(path)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	evs := l2.Events()
	if len(evs) != 1 {
		t.Fatalf("expected 1 event, got %d", len(evs))
	}
	if evs[0].Kind != audit.KindOpened {
		t.Errorf("kind = %s, want opened", evs[0].Kind)
	}
	if evs[0].Port.Port != 80 {
		t.Errorf("port = %d, want 80", evs[0].Port.Port)
	}
}

func TestLoadMissingFileReturnsEmpty(t *testing.T) {
	l, err := audit.New("/tmp/portwatch_audit_missing_xyz.json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(l.Events()) != 0 {
		t.Error("expected empty log")
	}
}

func TestLoadInvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	_ = os.WriteFile(path, []byte("not json"), 0o644)
	_, err := audit.New(path)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestMultipleAppends(t *testing.T) {
	dir := t.TempDir()
	l, _ := audit.New(filepath.Join(dir, "audit.json"))
	_ = l.Append(audit.KindOpened, makePort("tcp", 443), "")
	_ = l.Append(audit.KindClosed, makePort("tcp", 443), "gone")
	_ = l.Append(audit.KindAlert, makePort("udp", 53), "suspicious")
	evs := l.Events()
	if len(evs) != 3 {
		t.Fatalf("expected 3 events, got %d", len(evs))
	}
	if evs[2].Message != "suspicious" {
		t.Errorf("message = %q, want suspicious", evs[2].Message)
	}
}
