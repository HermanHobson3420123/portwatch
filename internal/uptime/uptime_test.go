package uptime_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"portwatch/internal/uptime"
)

func TestOpenedCreatesRecord(t *testing.T) {
	tr, _ := uptime.New(filepath.Join(t.TempDir(), "uptime.json"))
	now := time.Now()
	tr.Opened("tcp", 80, now)
	r, ok := tr.Get("tcp", 80)
	if !ok {
		t.Fatal("expected record")
	}
	if !r.FirstSeen.Equal(now) {
		t.Errorf("FirstSeen = %v, want %v", r.FirstSeen, now)
	}
}

func TestOpenedIdempotent(t *testing.T) {
	tr, _ := uptime.New(filepath.Join(t.TempDir(), "uptime.json"))
	t1 := time.Now()
	t2 := t1.Add(time.Minute)
	tr.Opened("tcp", 80, t1)
	tr.Opened("tcp", 80, t2) // should not overwrite
	r, _ := tr.Get("tcp", 80)
	if !r.FirstSeen.Equal(t1) {
		t.Errorf("FirstSeen overwritten: got %v, want %v", r.FirstSeen, t1)
	}
}

func TestSeenUpdatesLastSeen(t *testing.T) {
	tr, _ := uptime.New(filepath.Join(t.TempDir(), "uptime.json"))
	now := time.Now()
	tr.Opened("tcp", 443, now)
	later := now.Add(5 * time.Minute)
	tr.Seen("tcp", 443, later)
	r, _ := tr.Get("tcp", 443)
	if !r.LastSeen.Equal(later) {
		t.Errorf("LastSeen = %v, want %v", r.LastSeen, later)
	}
}

func TestClosedRemovesRecord(t *testing.T) {
	tr, _ := uptime.New(filepath.Join(t.TempDir(), "uptime.json"))
	tr.Opened("udp", 53, time.Now())
	tr.Closed("udp", 53)
	_, ok := tr.Get("udp", 53)
	if ok {
		t.Error("expected record to be removed")
	}
}

func TestDurationCalculation(t *testing.T) {
	tr, _ := uptime.New(filepath.Join(t.TempDir(), "uptime.json"))
	now := time.Now()
	tr.Opened("tcp", 22, now)
	tr.Seen("tcp", 22, now.Add(2*time.Hour))
	r, _ := tr.Get("tcp", 22)
	if r.Duration() != 2*time.Hour {
		t.Errorf("Duration = %v, want 2h", r.Duration())
	}
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "uptime.json")
	tr, _ := uptime.New(path)
	now := time.Now().Truncate(time.Second)
	tr.Opened("tcp", 8080, now)
	if err := tr.Save(); err != nil {
		t.Fatal(err)
	}
	tr2, err := uptime.New(path)
	if err != nil {
		t.Fatal(err)
	}
	r, ok := tr2.Get("tcp", 8080)
	if !ok {
		t.Fatal("record not persisted")
	}
	if !r.FirstSeen.Equal(now) {
		t.Errorf("FirstSeen = %v, want %v", r.FirstSeen, now)
	}
}

func TestLoadMissingFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing.json")
	tr, err := uptime.New(path)
	if err != nil {
		t.Fatal(err)
	}
	_ = tr
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Error("file should not be created on load")
	}
}
