package portcount

import (
	"strings"
	"testing"

	"portwatch/internal/scanner"
)

func ports(specs [][2]string) []scanner.Port {
	out := make([]scanner.Port, 0, len(specs))
	for i, s := range specs {
		out = append(out, scanner.Port{Port: i + 1, Protocol: s[0], Address: s[1]})
	}
	return out
}

func TestRecordCountsTotal(t *testing.T) {
	c := New()
	snap := c.Record(ports([][2]string{{"tcp", "127.0.0.1"}, {"tcp", "127.0.0.1"}, {"udp", "127.0.0.1"}}))
	if snap.Total != 3 {
		t.Fatalf("expected total 3, got %d", snap.Total)
	}
}

func TestRecordCountsTCP(t *testing.T) {
	c := New()
	snap := c.Record(ports([][2]string{{"tcp", "127.0.0.1"}, {"tcp", "127.0.0.1"}}))
	if snap.TCP != 2 {
		t.Fatalf("expected tcp 2, got %d", snap.TCP)
	}
	if snap.UDP != 0 {
		t.Fatalf("expected udp 0, got %d", snap.UDP)
	}
}

func TestRecordCountsUDP(t *testing.T) {
	c := New()
	snap := c.Record(ports([][2]string{{"udp", "127.0.0.1"}, {"udp", "127.0.0.1"}, {"udp", "127.0.0.1"}}))
	if snap.UDP != 3 {
		t.Fatalf("expected udp 3, got %d", snap.UDP)
	}
}

func TestLastReturnsLatest(t *testing.T) {
	c := New()
	c.Record(ports([][2]string{{"tcp", "127.0.0.1"}}))
	c.Record(ports([][2]string{{"tcp", "127.0.0.1"}, {"udp", "127.0.0.1"}}))
	s := c.Last()
	if s.Total != 2 {
		t.Fatalf("expected last total 2, got %d", s.Total)
	}
}

func TestSummaryContainsFields(t *testing.T) {
	c := New()
	c.Record(ports([][2]string{{"tcp", "127.0.0.1"}, {"udp", "127.0.0.1"}}))
	summ := c.Summary()
	for _, want := range []string{"total=2", "tcp=1", "udp=1"} {
		if !strings.Contains(summ, want) {
			t.Errorf("summary missing %q: %s", want, summ)
		}
	}
}

func TestEmptyPortsReturnsZeroSnapshot(t *testing.T) {
	c := New()
	snap := c.Record(nil)
	if snap.Total != 0 || snap.TCP != 0 || snap.UDP != 0 {
		t.Fatalf("expected zero snapshot, got %+v", snap)
	}
}
