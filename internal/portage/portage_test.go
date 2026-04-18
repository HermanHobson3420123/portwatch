package portage

import (
	"testing"
	"time"

	"portwatch/internal/scanner"
)

func p(proto string, port int) scanner.Port {
	return scanner.Port{Proto: proto, Port: port}
}

func TestOpenedCreatesEntry(t *testing.T) {
	tr := New()
	tr.Opened(p("tcp", 80))
	e, ok := tr.Get(p("tcp", 80))
	if !ok {
		t.Fatal("expected entry")
	}
	if e.Port.Port != 80 {
		t.Errorf("expected port 80, got %d", e.Port.Port)
	}
}

func TestOpenedIdempotent(t *testing.T) {
	tr := New()
	fixed := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	tr.now = func() time.Time { return fixed }
	tr.Opened(p("tcp", 443))
	tr.now = func() time.Time { return fixed.Add(time.Hour) }
	tr.Opened(p("tcp", 443))
	e, _ := tr.Get(p("tcp", 443))
	if !e.FirstSeen.Equal(fixed) {
		t.Error("second Opened should not overwrite first-seen time")
	}
}

func TestClosedRemovesEntry(t *testing.T) {
	tr := New()
	tr.Opened(p("tcp", 22))
	tr.Closed(p("tcp", 22))
	_, ok := tr.Get(p("tcp", 22))
	if ok {
		t.Error("expected entry to be removed")
	}
}

func TestAgeString(t *testing.T) {
	now := time.Now()
	cases := []struct {
		dur time.Duration
		want string
	}{
		{30 * time.Second, "30s"},
		{90 * time.Second, "1m"},
		{2 * time.Hour, "2h"},
	}
	for _, c := range cases {
		e := Entry{FirstSeen: now.Add(-c.dur)}
		got := e.AgeString(now)
		if got != c.want {
			t.Errorf("AgeString(%v) = %q, want %q", c.dur, got, c.want)
		}
	}
}

func TestAllReturnsEntries(t *testing.T) {
	tr := New()
	tr.Opened(p("tcp", 80))
	tr.Opened(p("udp", 53))
	all := tr.All()
	if len(all) != 2 {
		t.Errorf("expected 2 entries, got %d", len(all))
	}
}
