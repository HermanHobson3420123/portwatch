package metrics

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRecordScanIncrements(t *testing.T) {
	c := New()
	c.RecordScan()
	c.RecordScan()
	if got := c.Snapshot().ScansTotal; got != 2 {
		t.Fatalf("expected 2 scans, got %d", got)
	}
}

func TestRecordOpenedAndClosed(t *testing.T) {
	c := New()
	c.RecordOpened()
	c.RecordClosed()
	c.RecordClosed()
	s := c.Snapshot()
	if s.OpenedTotal != 1 {
		t.Fatalf("expected 1 opened, got %d", s.OpenedTotal)
	}
	if s.ClosedTotal != 2 {
		t.Fatalf("expected 2 closed, got %d", s.ClosedTotal)
	}
	if s.AlertsTotal != 3 {
		t.Fatalf("expected 3 alerts, got %d", s.AlertsTotal)
	}
}

func TestLastScanUpdated(t *testing.T) {
	c := New()
	if !c.Snapshot().LastScan.IsZero() {
		t.Fatal("expected zero LastScan before any scan")
	}
	c.RecordScan()
	if c.Snapshot().LastScan.IsZero() {
		t.Fatal("expected non-zero LastScan after scan")
	}
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "metrics.json")

	c := New()
	c.RecordScan()
	c.RecordOpened()

	if err := c.Save(path); err != nil {
		t.Fatalf("Save: %v", err)
	}

	s, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if s.ScansTotal != 1 || s.OpenedTotal != 1 {
		t.Fatalf("unexpected snapshot: %+v", s)
	}
}

func TestLoadMissingFile(t *testing.T) {
	_, err := Load(filepath.Join(t.TempDir(), "nope.json"))
	if !os.IsNotExist(err) {
		t.Fatalf("expected not-exist error, got %v", err)
	}
}
