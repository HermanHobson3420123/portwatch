package healthcheck

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "health.json")
}

func TestNewDefaults(t *testing.T) {
	c := New(tempPath(t))
	s := c.Get()
	if !s.Healthy {
		t.Fatal("expected healthy by default")
	}
	if s.ScanCount != 0 || s.AlertCount != 0 {
		t.Fatal("expected zero counters")
	}
}

func TestRecordScanIncrementsCount(t *testing.T) {
	c := New(tempPath(t))
	c.RecordScan()
	c.RecordScan()
	if got := c.Get().ScanCount; got != 2 {
		t.Fatalf("expected 2, got %d", got)
	}
}

func TestRecordAlertIncrementsCount(t *testing.T) {
	c := New(tempPath(t))
	c.RecordAlert()
	if got := c.Get().AlertCount; got != 1 {
		t.Fatalf("expected 1, got %d", got)
	}
}

func TestSetHealthy(t *testing.T) {
	c := New(tempPath(t))
	c.SetHealthy(false)
	if c.Get().Healthy {
		t.Fatal("expected unhealthy")
	}
	c.SetHealthy(true)
	if !c.Get().Healthy {
		t.Fatal("expected healthy")
	}
}

func TestFlushWritesJSON(t *testing.T) {
	p := tempPath(t)
	c := New(p)
	c.RecordScan()
	c.RecordAlert()

	data, err := os.ReadFile(p)
	if err != nil {
		t.Fatalf("read file: %v", err)
	}
	var s Status
	if err := json.Unmarshal(data, &s); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if s.ScanCount != 1 || s.AlertCount != 1 {
		t.Fatalf("unexpected counts: %+v", s)
	}
}

func TestLastScanUpdated(t *testing.T) {
	c := New(tempPath(t))
	before := time.Now()
	c.RecordScan()
	after := time.Now()
	ls := c.Get().LastScan
	if ls.Before(before) || ls.After(after) {
		t.Fatalf("last scan timestamp out of range: %v", ls)
	}
}
