package portstate

import (
	"testing"
	"time"

	"portwatch/internal/scanner"
)

func p(proto string, port int) scanner.Port {
	return scanner.Port{Proto: proto, Number: port}
}

var t0 = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func TestOpenCreatesEntry(t *testing.T) {
	tr := New()
	tr.Open(p("tcp", 80), t0)
	if tr.Len() != 1 {
		t.Fatalf("expected 1 entry, got %d", tr.Len())
	}
	s, ok := tr.Get(p("tcp", 80))
	if !ok {
		t.Fatal("expected entry to exist")
	}
	if !s.OpenSince.Equal(t0) {
		t.Errorf("unexpected OpenSince: %v", s.OpenSince)
	}
}

func TestOpenIdempotent(t *testing.T) {
	tr := New()
	tr.Open(p("tcp", 80), t0)
	tr.Open(p("tcp", 80), t0.Add(time.Minute))
	if tr.Len() != 1 {
		t.Fatalf("expected 1 entry, got %d", tr.Len())
	}
	s, _ := tr.Get(p("tcp", 80))
	if !s.LastSeen.Equal(t0.Add(time.Minute)) {
		t.Errorf("LastSeen not updated: %v", s.LastSeen)
	}
	if !s.OpenSince.Equal(t0) {
		t.Errorf("OpenSince should not change: %v", s.OpenSince)
	}
}

func TestCloseRemovesEntry(t *testing.T) {
	tr := New()
	tr.Open(p("tcp", 443), t0)
	tr.Close(p("tcp", 443))
	if tr.Len() != 0 {
		t.Fatal("expected 0 entries after close")
	}
	_, ok := tr.Get(p("tcp", 443))
	if ok {
		t.Fatal("entry should not exist after close")
	}
}

func TestGetMissingReturnsFalse(t *testing.T) {
	tr := New()
	_, ok := tr.Get(p("udp", 53))
	if ok {
		t.Fatal("expected false for missing port")
	}
}

func TestAllReturnsSnapshot(t *testing.T) {
	tr := New()
	tr.Open(p("tcp", 22), t0)
	tr.Open(p("tcp", 80), t0)
	tr.Open(p("udp", 53), t0)
	all := tr.All()
	if len(all) != 3 {
		t.Fatalf("expected 3 states, got %d", len(all))
	}
}

func TestProtocolsAreDistinct(t *testing.T) {
	tr := New()
	tr.Open(p("tcp", 53), t0)
	tr.Open(p("udp", 53), t0)
	if tr.Len() != 2 {
		t.Fatalf("tcp and udp on same port should be distinct, got %d", tr.Len())
	}
}
