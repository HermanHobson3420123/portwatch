package portclassify_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/portclassify"
	"github.com/user/portwatch/internal/scanner"
)

func p(number uint16, proto string) scanner.Port {
	return scanner.Port{Number: number, Protocol: proto}
}

func TestClassifySystemPort(t *testing.T) {
	c := portclassify.Classify(p(22, "tcp"))
	if c.Tier != portclassify.TierSystem {
		t.Fatalf("expected system tier, got %s", c.Tier)
	}
	if c.Service != "ssh" {
		t.Fatalf("expected ssh service, got %q", c.Service)
	}
}

func TestClassifyRegisteredPort(t *testing.T) {
	c := portclassify.Classify(p(3306, "tcp"))
	if c.Tier != portclassify.TierRegistered {
		t.Fatalf("expected registered tier, got %s", c.Tier)
	}
	if c.Service != "mysql" {
		t.Fatalf("expected mysql, got %q", c.Service)
	}
}

func TestClassifyDynamicPort(t *testing.T) {
	c := portclassify.Classify(p(55000, "tcp"))
	if c.Tier != portclassify.TierDynamic {
		t.Fatalf("expected dynamic tier, got %s", c.Tier)
	}
	if c.Service != "" {
		t.Fatalf("expected empty service for dynamic port, got %q", c.Service)
	}
}

func TestClassifyUnknownService(t *testing.T) {
	c := portclassify.Classify(p(9999, "tcp"))
	if c.Tier != portclassify.TierRegistered {
		t.Fatalf("expected registered tier, got %s", c.Tier)
	}
	if c.Service != "" {
		t.Fatalf("expected empty service, got %q", c.Service)
	}
}

func TestClassifyAll(t *testing.T) {
	ports := []scanner.Port{p(80, "tcp"), p(443, "tcp"), p(60000, "udp")}
	classes := portclassify.ClassifyAll(ports)
	if len(classes) != 3 {
		t.Fatalf("expected 3 classes, got %d", len(classes))
	}
}

func TestSummary(t *testing.T) {
	classes := portclassify.ClassifyAll([]scanner.Port{
		p(22, "tcp"), p(80, "tcp"), p(3306, "tcp"), p(55000, "tcp"),
	})
	s := portclassify.Summary(classes)
	if !strings.Contains(s, "system=2") {
		t.Errorf("expected system=2 in summary, got %q", s)
	}
	if !strings.Contains(s, "registered=1") {
		t.Errorf("expected registered=1 in summary, got %q", s)
	}
	if !strings.Contains(s, "dynamic=1") {
		t.Errorf("expected dynamic=1 in summary, got %q", s)
	}
}

func TestPrintReportContainsHeaders(t *testing.T) {
	var buf bytes.Buffer
	classes := portclassify.ClassifyAll([]scanner.Port{p(22, "tcp")})
	portclassify.PrintReport(&buf, classes)
	out := buf.String()
	for _, hdr := range []string{"PORT", "PROTO", "TIER", "SERVICE"} {
		if !strings.Contains(out, hdr) {
			t.Errorf("expected header %q in output", hdr)
		}
	}
}

func TestPrintReportEmpty(t *testing.T) {
	var buf bytes.Buffer
	portclassify.PrintReport(&buf, nil)
	if !strings.Contains(buf.String(), "no ports") {
		t.Errorf("expected 'no ports' message for empty input")
	}
}
