package topports

import (
	"testing"

	"portwatch/internal/scanner"
)

func p(number int, proto string) scanner.Port {
	return scanner.Port{Number: number, Protocol: proto}
}

func TestRecordAndTopBasic(t *testing.T) {
	c := New()
	c.Record([]scanner.Port{p(80, "tcp"), p(443, "tcp")})
	c.Record([]scanner.Port{p(80, "tcp")})

	top := c.Top(0)
	if len(top) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(top))
	}
	if top[0].Port.Number != 80 || top[0].Count != 2 {
		t.Errorf("expected port 80 with count 2, got %+v", top[0])
	}
	if top[1].Port.Number != 443 || top[1].Count != 1 {
		t.Errorf("expected port 443 with count 1, got %+v", top[1])
	}
}

func TestTopLimitsResults(t *testing.T) {
	c := New()
	for _, port := range []int{22, 80, 443, 8080, 9090} {
		c.Record([]scanner.Port{p(port, "tcp")})
	}
	top := c.Top(3)
	if len(top) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(top))
	}
}

func TestTopTieBreakByPortNumber(t *testing.T) {
	c := New()
	c.Record([]scanner.Port{p(443, "tcp"), p(80, "tcp")})

	top := c.Top(0)
	if top[0].Port.Number != 80 {
		t.Errorf("expected port 80 first on tie, got %d", top[0].Port.Number)
	}
}

func TestReset(t *testing.T) {
	c := New()
	c.Record([]scanner.Port{p(80, "tcp")})
	c.Reset()
	if len(c.Top(0)) != 0 {
		t.Error("expected empty after reset")
	}
}

func TestProtocolsTrackedSeparately(t *testing.T) {
	c := New()
	c.Record([]scanner.Port{p(53, "tcp"), p(53, "udp")})
	c.Record([]scanner.Port{p(53, "udp")})

	top := c.Top(0)
	if len(top) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(top))
	}
	var udpCount int
	for _, e := range top {
		if e.Port.Protocol == "udp" {
			udpCount = e.Count
		}
	}
	if udpCount != 2 {
		t.Errorf("expected udp count 2, got %d", udpCount)
	}
}
