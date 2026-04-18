package portfreq

import (
	"testing"
	"time"

	"portwatch/internal/scanner"
)

func p(proto, addr string) scanner.Port { return scanner.Port{Proto: proto, Addr: addr} }

var now = time.Now()

func TestRecordIncrementsCount(t *testing.T) {
	tr := New()
	tr.Record(p("tcp", "0.0.0.0:80"), now)
	tr.Record(p("tcp", "0.0.0.0:80"), now.Add(time.Second))
	top := tr.Top(10)
	if len(top) != 1 || top[0].SeenCount != 2 {
		t.Fatalf("expected count 2, got %+v", top)
	}
}

func TestTopLimitsResults(t *testing.T) {
	tr := New()
	for i := 0; i < 5; i++ {
		tr.Record(p("tcp", "0.0.0.0:80"), now)
	}
	tr.Record(p("tcp", "0.0.0.0:443"), now)
	top := tr.Top(1)
	if len(top) != 1 {
		t.Fatalf("expected 1 result, got %d", len(top))
	}
	if top[0].Port.Addr != "0.0.0.0:80" {
		t.Fatalf("expected port 80 at top, got %s", top[0].Port.Addr)
	}
}

func TestTopSortedByCount(t *testing.T) {
	tr := New()
	tr.Record(p("tcp", "0.0.0.0:22"), now)
	for i := 0; i < 3; i++ {
		tr.Record(p("tcp", "0.0.0.0:443"), now)
	}
	top := tr.Top(0)
	if top[0].Port.Addr != "0.0.0.0:443" {
		t.Fatalf("expected 443 first, got %s", top[0].Port.Addr)
	}
}

func TestResetClearsData(t *testing.T) {
	tr := New()
	tr.Record(p("tcp", "0.0.0.0:80"), now)
	tr.Reset()
	if len(tr.Top(10)) != 0 {
		t.Fatal("expected empty after reset")
	}
}

func TestFirstAndLastSeen(t *testing.T) {
	tr := New()
	t1 := now
	t2 := now.Add(time.Minute)
	tr.Record(p("tcp", "0.0.0.0:80"), t1)
	tr.Record(p("tcp", "0.0.0.0:80"), t2)
	e := tr.Top(1)[0]
	if !e.FirstSeen.Equal(t1) {
		t.Fatalf("expected FirstSeen %v, got %v", t1, e.FirstSeen)
	}
	if !e.LastSeen.Equal(t2) {
		t.Fatalf("expected LastSeen %v, got %v", t2, e.LastSeen)
	}
}
