package portschedule

import (
	"testing"
	"time"
)

var base = time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

func makeEntry(label string, offsetStart, offsetEnd time.Duration) Entry {
	return Entry{
		Label:    label,
		Start:    base.Add(offsetStart),
		End:      base.Add(offsetEnd),
		PortLow:  1,
		PortHigh: 1024,
		Protocol: "tcp",
	}
}

func TestAddAndAll(t *testing.T) {
	s := New()
	s.Add(makeEntry("a", -time.Hour, time.Hour))
	s.Add(makeEntry("b", time.Hour, 2*time.Hour))
	if len(s.All()) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(s.All()))
	}
}

func TestActiveNow(t *testing.T) {
	s := New()
	s.Add(makeEntry("active", -time.Hour, time.Hour))
	s.Add(makeEntry("future", time.Hour, 2*time.Hour))
	s.Add(makeEntry("past", -2*time.Hour, -time.Hour))

	active := s.ActiveNow(base)
	if len(active) != 1 {
		t.Fatalf("expected 1 active entry, got %d", len(active))
	}
	if active[0].Label != "active" {
		t.Errorf("expected 'active', got %q", active[0].Label)
	}
}

func TestRemoveByLabel(t *testing.T) {
	s := New()
	s.Add(makeEntry("keep", -time.Hour, time.Hour))
	s.Add(makeEntry("drop", -time.Hour, time.Hour))
	s.Add(makeEntry("drop", time.Hour, 2*time.Hour))
	s.Remove("drop")
	all := s.All()
	if len(all) != 1 || all[0].Label != "keep" {
		t.Errorf("expected only 'keep' entry, got %+v", all)
	}
}

func TestActiveEntryBoundary(t *testing.T) {
	e := makeEntry("boundary", 0, time.Hour)
	if !e.Active(base) {
		t.Error("expected entry active at start boundary")
	}
	if e.Active(base.Add(time.Hour)) {
		t.Error("expected entry inactive at end boundary")
	}
}

func TestActiveNowEmpty(t *testing.T) {
	s := New()
	active := s.ActiveNow(base)
	if len(active) != 0 {
		t.Errorf("expected no active entries on empty schedule, got %d", len(active))
	}
}
