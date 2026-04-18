package portevict

import (
	"bytes"
	"testing"
	"time"

	"portwatch/internal/scanner"
)

func p(num int, proto string) scanner.Port {
	return scanner.Port{Number: num, Protocol: proto}
}

func TestClosedBelowThresholdRecordsEviction(t *testing.T) {
	tr := New(5 * time.Second)
	now := time.Now()
	evicted := tr.Closed(p(8080, "tcp"), now, now.Add(2*time.Second))
	if !evicted {
		t.Fatal("expected eviction")
	}
	if got := len(tr.All()); got != 1 {
		t.Fatalf("expected 1 entry, got %d", got)
	}
}

func TestClosedAboveThresholdNotEvicted(t *testing.T) {
	tr := New(5 * time.Second)
	now := time.Now()
	evicted := tr.Closed(p(443, "tcp"), now, now.Add(10*time.Second))
	if evicted {
		t.Fatal("should not be evicted")
	}
	if got := len(tr.All()); got != 0 {
		t.Fatalf("expected 0 entries, got %d", got)
	}
}

func TestResetClearsEntries(t *testing.T) {
	tr := New(5 * time.Second)
	now := time.Now()
	tr.Closed(p(22, "tcp"), now, now.Add(1*time.Second))
	tr.Reset()
	if got := len(tr.All()); got != 0 {
		t.Fatalf("expected 0 after reset, got %d", got)
	}
}

func TestAllReturnsCopy(t *testing.T) {
	tr := New(time.Minute)
	now := time.Now()
	tr.Closed(p(80, "tcp"), now, now.Add(500*time.Millisecond))
	a := tr.All()
	a[0].Port.Number = 9999
	if tr.All()[0].Port.Number == 9999 {
		t.Fatal("All() should return a copy")
	}
}

func TestPrintEvictionsNoEntries(t *testing.T) {
	var buf bytes.Buffer
	PrintEvictions(&buf, nil)
	if buf.Len() == 0 {
		t.Fatal("expected output")
	}
}

func TestSummaryNonEmpty(t *testing.T) {
	tr := New(time.Minute)
	now := time.Now()
	tr.Closed(p(3000, "tcp"), now, now.Add(100*time.Millisecond))
	s := Summary(tr.All())
	if s == "evictions: none" {
		t.Fatal("expected non-empty summary")
	}
}
