package portwindow

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"portwatch/internal/scanner"
)

func p(proto, addr string) scanner.Port {
	return scanner.Port{Proto: proto, Addr: addr}
}

func fixedNow(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestRecordAndActive(t *testing.T) {
	w := New(5 * time.Minute)
	w.Record(p("tcp", "0.0.0.0:80"))
	w.Record(p("tcp", "0.0.0.0:443"))
	if got := w.Len(); got != 2 {
		t.Fatalf("expected 2 entries, got %d", got)
	}
}

func TestRecordIncrementsCount(t *testing.T) {
	w := New(5 * time.Minute)
	w.Record(p("tcp", "0.0.0.0:8080"))
	w.Record(p("tcp", "0.0.0.0:8080"))
	entries := w.Active()
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Count != 2 {
		t.Errorf("expected count 2, got %d", entries[0].Count)
	}
}

func TestEvictExpiredEntries(t *testing.T) {
	base := time.Now()
	w := New(1 * time.Minute)
	w.now = fixedNow(base)
	w.Record(p("tcp", "0.0.0.0:22"))
	// advance past the window
	w.now = fixedNow(base.Add(2 * time.Minute))
	if got := w.Len(); got != 0 {
		t.Errorf("expected 0 after eviction, got %d", got)
	}
}

func TestActiveReturnsLiveEntries(t *testing.T) {
	base := time.Now()
	w := New(5 * time.Minute)
	w.now = fixedNow(base)
	w.Record(p("udp", "0.0.0.0:53"))
	w.now = fixedNow(base.Add(3 * time.Minute))
	w.Record(p("tcp", "0.0.0.0:80"))
	// advance 4 min from base: udp:53 last seen at base (4 min ago > 5 min? no, 4 < 5)
	w.now = fixedNow(base.Add(4 * time.Minute))
	if got := w.Len(); got != 2 {
		t.Errorf("expected 2, got %d", got)
	}
}

func TestPrintActiveEmpty(t *testing.T) {
	var buf bytes.Buffer
	PrintActive(&buf, nil)
	if !strings.Contains(buf.String(), "no ports") {
		t.Errorf("unexpected output: %q", buf.String())
	}
}

func TestSummaryNonEmpty(t *testing.T) {
	w := New(5 * time.Minute)
	w.Record(p("tcp", "0.0.0.0:80"))
	w.Record(p("tcp", "0.0.0.0:80"))
	s := Summary(w.Active())
	if !strings.Contains(s, "1 active port") {
		t.Errorf("unexpected summary: %q", s)
	}
	if !strings.Contains(s, "2 total") {
		t.Errorf("expected total observations in summary: %q", s)
	}
}

func TestSummaryEmpty(t *testing.T) {
	s := Summary(nil)
	if !strings.Contains(s, "0 active") {
		t.Errorf("unexpected summary: %q", s)
	}
}
