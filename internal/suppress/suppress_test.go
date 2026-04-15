package suppress

import (
	"testing"
	"time"
)

func fixedNow(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestIsSuppressed_ActiveEntry(t *testing.T) {
	l := New()
	now := time.Now()
	l.now = fixedNow(now)
	l.Add(8080, "tcp", now.Add(10*time.Minute))

	if !l.IsSuppressed(8080, "tcp") {
		t.Error("expected port 8080/tcp to be suppressed")
	}
}

func TestIsSuppressed_ExpiredEntry(t *testing.T) {
	l := New()
	now := time.Now()
	l.now = fixedNow(now)
	l.Add(8080, "tcp", now.Add(-1*time.Minute))

	if l.IsSuppressed(8080, "tcp") {
		t.Error("expected expired entry to not be suppressed")
	}
}

func TestIsSuppressed_MissingEntry(t *testing.T) {
	l := New()
	if l.IsSuppressed(443, "tcp") {
		t.Error("expected unknown port to not be suppressed")
	}
}

func TestRemove(t *testing.T) {
	l := New()
	now := time.Now()
	l.now = fixedNow(now)
	l.Add(80, "tcp", now.Add(1*time.Hour))
	l.Remove(80, "tcp")

	if l.IsSuppressed(80, "tcp") {
		t.Error("expected removed entry to not be suppressed")
	}
}

func TestPurgeRemovesExpired(t *testing.T) {
	l := New()
	now := time.Now()
	l.now = fixedNow(now)
	l.Add(80, "tcp", now.Add(-1*time.Minute))
	l.Add(443, "tcp", now.Add(1*time.Hour))
	l.Purge()

	if l.IsSuppressed(80, "tcp") {
		t.Error("expected expired entry to be purged")
	}
	if !l.IsSuppressed(443, "tcp") {
		t.Error("expected active entry to remain after purge")
	}
}

func TestActiveReturnsOnlyLiveEntries(t *testing.T) {
	l := New()
	now := time.Now()
	l.now = fixedNow(now)
	l.Add(80, "tcp", now.Add(-1*time.Minute))
	l.Add(443, "tcp", now.Add(1*time.Hour))
	l.Add(22, "tcp", now.Add(30*time.Minute))

	active := l.Active()
	if len(active) != 2 {
		t.Errorf("expected 2 active entries, got %d", len(active))
	}
}

func TestProtocolIndependence(t *testing.T) {
	l := New()
	now := time.Now()
	l.now = fixedNow(now)
	l.Add(53, "tcp", now.Add(1*time.Hour))

	if l.IsSuppressed(53, "udp") {
		t.Error("tcp suppression should not affect udp")
	}
}
