package portannounce_test

import (
	"strings"
	"testing"

	"github.com/user/portwatch/internal/portannounce"
	"github.com/user/portwatch/internal/scanner"
)

func p(proto string, num int) scanner.Port {
	return scanner.Port{Protocol: proto, Number: num}
}

func TestAnnounceWritesHeader(t *testing.T) {
	var buf strings.Builder
	a := portannounce.New(&buf)
	a.Announce([]scanner.Port{p("tcp", 80)})
	if !strings.Contains(buf.String(), "portwatch startup") {
		t.Errorf("expected startup header, got: %s", buf.String())
	}
}

func TestAnnounceListsPorts(t *testing.T) {
	var buf strings.Builder
	a := portannounce.New(&buf)
	a.Announce([]scanner.Port{p("tcp", 443), p("udp", 53)})
	out := buf.String()
	if !strings.Contains(out, "TCP") || !strings.Contains(out, "443") {
		t.Errorf("expected TCP 443 in output, got: %s", out)
	}
	if !strings.Contains(out, "UDP") || !strings.Contains(out, "53") {
		t.Errorf("expected UDP 53 in output, got: %s", out)
	}
}

func TestAnnounceEmptyPorts(t *testing.T) {
	var buf strings.Builder
	a := portannounce.New(&buf)
	a.Announce(nil)
	out := buf.String()
	if !strings.Contains(out, "(none)") {
		t.Errorf("expected (none) for empty port list, got: %s", out)
	}
}

func TestAnnounceSortsByProtocolThenNumber(t *testing.T) {
	var buf strings.Builder
	a := portannounce.New(&buf)
	summary := a.Announce([]scanner.Port{p("udp", 53), p("tcp", 8080), p("tcp", 80)})
	if len(summary.Ports) != 3 {
		t.Fatalf("expected 3 ports, got %d", len(summary.Ports))
	}
	if summary.Ports[0].Protocol != "tcp" || summary.Ports[0].Number != 80 {
		t.Errorf("first port should be tcp/80, got %+v", summary.Ports[0])
	}
	if summary.Ports[1].Protocol != "tcp" || summary.Ports[1].Number != 8080 {
		t.Errorf("second port should be tcp/8080, got %+v", summary.Ports[1])
	}
	if summary.Ports[2].Protocol != "udp" || summary.Ports[2].Number != 53 {
		t.Errorf("third port should be udp/53, got %+v", summary.Ports[2])
	}
}

func TestAnnounceSummaryTimestampSet(t *testing.T) {
	var buf strings.Builder
	a := portannounce.New(&buf)
	summary := a.Announce([]scanner.Port{p("tcp", 22)})
	if summary.At.IsZero() {
		t.Error("expected non-zero timestamp in summary")
	}
}
