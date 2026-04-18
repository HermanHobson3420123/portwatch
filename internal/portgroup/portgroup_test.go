package portgroup

import (
	"strings"
	"testing"

	"portwatch/internal/scanner"
)

func p(num int, proto string) scanner.Port {
	return scanner.Port{Number: num, Protocol: proto}
}

func TestByProtocolGroupsTCP(t *testing.T) {
	ports := []scanner.Port{p(80, "tcp"), p(443, "tcp"), p(53, "udp")}
	groups := ByProtocol().Group(ports)
	if len(groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(groups))
	}
	if groups[0].Name != "tcp" {
		t.Errorf("expected first group tcp, got %s", groups[0].Name)
	}
	if len(groups[0].Ports) != 2 {
		t.Errorf("expected 2 tcp ports, got %d", len(groups[0].Ports))
	}
}

func TestByPortRangeBuckets(t *testing.T) {
	ports := []scanner.Port{p(22, "tcp"), p(8080, "tcp"), p(51000, "udp")}
	groups := ByPortRange().Group(ports)
	names := make(map[string]int)
	for _, g := range groups {
		names[g.Name] = len(g.Ports)
	}
	if names["system (0-1023)"] != 1 {
		t.Errorf("expected 1 system port")
	}
	if names["registered (1024-49151)"] != 1 {
		t.Errorf("expected 1 registered port")
	}
	if names["dynamic (49152+)"] != 1 {
		t.Errorf("expected 1 dynamic port")
	}
}

func TestGroupSortedByPortNumber(t *testing.T) {
	ports := []scanner.Port{p(443, "tcp"), p(80, "tcp"), p(22, "tcp")}
	groups := ByProtocol().Group(ports)
	ps := groups[0].Ports
	if ps[0].Number != 22 || ps[1].Number != 80 || ps[2].Number != 443 {
		t.Errorf("ports not sorted: %v", ps)
	}
}

func TestSummaryContainsGroupNames(t *testing.T) {
	ports := []scanner.Port{p(80, "tcp"), p(53, "udp")}
	groups := ByProtocol().Group(ports)
	s := Summary(groups)
	if !strings.Contains(s, "tcp") || !strings.Contains(s, "udp") {
		t.Errorf("summary missing group names: %s", s)
	}
}

func TestEmptyPortsReturnsNoGroups(t *testing.T) {
	groups := ByProtocol().Group(nil)
	if len(groups) != 0 {
		t.Errorf("expected no groups, got %d", len(groups))
	}
}
