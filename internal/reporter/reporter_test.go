package reporter_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/reporter"
	"github.com/user/portwatch/internal/scanner"
)

func makePort(proto string, num int) scanner.Port {
	return scanner.Port{Protocol: proto, Number: num}
}

func TestReportOpenedText(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.New(&buf, reporter.FormatText)

	if err := r.ReportOpened(makePort("tcp", 8080)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "[opened]") {
		t.Errorf("expected '[opened]' in output, got: %s", out)
	}
	if !strings.Contains(out, "+") {
		t.Errorf("expected '+' symbol in output, got: %s", out)
	}
	if !strings.Contains(out, "8080") {
		t.Errorf("expected port number in output, got: %s", out)
	}
}

func TestReportClosedText(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.New(&buf, reporter.FormatText)

	if err := r.ReportClosed(makePort("udp", 53)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "[closed]") {
		t.Errorf("expected '[closed]' in output, got: %s", out)
	}
	if !strings.Contains(out, "-") {
		t.Errorf("expected '-' symbol in output, got: %s", out)
	}
}

func TestReportOpenedJSON(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.New(&buf, reporter.FormatJSON)

	if err := r.ReportOpened(makePort("tcp", 443)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var event reporter.Event
	if err := json.Unmarshal(buf.Bytes(), &event); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if event.Kind != "opened" {
		t.Errorf("expected kind=opened, got %s", event.Kind)
	}
	if event.Port.Number != 443 {
		t.Errorf("expected port 443, got %d", event.Port.Number)
	}
}

func TestReportClosedJSON(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.New(&buf, reporter.FormatJSON)

	if err := r.ReportClosed(makePort("tcp", 22)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var event reporter.Event
	if err := json.Unmarshal(buf.Bytes(), &event); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if event.Kind != "closed" {
		t.Errorf("expected kind=closed, got %s", event.Kind)
	}
	if event.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}
