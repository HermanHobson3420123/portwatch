package portexpiry

import (
	"testing"
	"time"

	"portwatch/internal/scanner"
)

func p(proto string, port int) scanner.Port {
	return scanner.Port{Proto: proto, Number: port}
}

var epoch = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func TestSeenPreventsExpiry(t *testing.T) {
	tr := New(5 * time.Second)
	tr.Seen(p("tcp", 80), epoch)
	expired := tr.Expired(epoch.Add(3 * time.Second))
	if len(expired) != 0 {
		t.Fatalf("expected no expired entries, got %d", len(expired))
	}
}

func TestExpiredAfterTTL(t *testing.T) {
	tr := New(5 * time.Second)
	tr.Seen(p("tcp", 80), epoch)
	expired := tr.Expired(epoch.Add(6 * time.Second))
	if len(expired) != 1 {
		t.Fatalf("expected 1 expired entry, got %d", len(expired))
	}
	if expired[0].Port.Number != 80 {
		t.Errorf("unexpected port %d", expired[0].Port.Number)
	}
}

func TestEvictRemovesEntries(t *testing.T) {
	tr := New(5 * time.Second)
	tr.Seen(p("tcp", 80), epoch)
	tr.Seen(p("tcp", 443), epoch)
	evicted := tr.Evict(epoch.Add(10 * time.Second))
	if len(evicted) != 2 {
		t.Fatalf("expected 2 evicted, got %d", len(evicted))
	}
	if tr.Len() != 0 {
		t.Errorf("expected tracker to be empty after evict")
	}
}

func TestEvictKeepsLiveEntries(t *testing.T) {
	tr := New(10 * time.Second)
	tr.Seen(p("tcp", 80), epoch)
	tr.Seen(p("tcp", 443), epoch.Add(8*time.Second))
	evicted := tr.Evict(epoch.Add(11 * time.Second))
	if len(evicted) != 1 {
		t.Fatalf("expected 1 evicted, got %d", len(evicted))
	}
	if tr.Len() != 1 {
		t.Errorf("expected 1 live entry remaining")
	}
}

func TestSeenUpdatesTimestamp(t *testing.T) {
	tr := New(5 * time.Second)
	tr.Seen(p("tcp", 80), epoch)
	tr.Seen(p("tcp", 80), epoch.Add(4*time.Second))
	expired := tr.Expired(epoch.Add(8 * time.Second))
	if len(expired) != 0 {
		t.Errorf("expected no expired after refresh, got %d", len(expired))
	}
}
