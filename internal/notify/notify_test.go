package notify_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"portwatch/internal/notify"
)

func TestSendDefaultTemplate(t *testing.T) {
	var buf bytes.Buffer
	n, err := notify.New(&buf)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	e := notify.Event{
		Timestamp: time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
		Level:     notify.LevelAlert,
		Message:   "port opened unexpectedly",
		Port:      8080,
		Protocol:  "tcp",
	}
	if err := n.Send(e); err != nil {
		t.Fatalf("Send: %v", err)
	}
	got := buf.String()
	if !strings.Contains(got, "ALERT") {
		t.Errorf("expected ALERT in output, got: %q", got)
	}
	if !strings.Contains(got, "8080/tcp") {
		t.Errorf("expected 8080/tcp in output, got: %q", got)
	}
	if !strings.Contains(got, "2024-06-01 12:00:00") {
		t.Errorf("expected timestamp in output, got: %q", got)
	}
}

func TestSendFillsTimestampWhenZero(t *testing.T) {
	var buf bytes.Buffer
	n, _ := notify.New(&buf)
	e := notify.Event{Level: notify.LevelInfo, Port: 22, Protocol: "tcp", Message: "ok"}
	if err := n.Send(e); err != nil {
		t.Fatalf("Send: %v", err)
	}
	if buf.Len() == 0 {
		t.Error("expected non-empty output")
	}
}

func TestNewWithTemplate(t *testing.T) {
	var buf bytes.Buffer
	n, err := notify.NewWithTemplate(&buf, "{{.Level}}:{{.Port}}\n")
	if err != nil {
		t.Fatalf("NewWithTemplate: %v", err)
	}
	_ = n.Send(notify.Event{Level: notify.LevelWarn, Port: 443, Protocol: "tcp"})
	got := strings.TrimSpace(buf.String())
	if got != "WARN:443" {
		t.Errorf("unexpected output: %q", got)
	}
}

func TestNewWithInvalidTemplate(t *testing.T) {
	_, err := notify.NewWithTemplate(nil, "{{.Unclosed")
	if err == nil {
		t.Error("expected error for invalid template")
	}
}

func TestWatchConsumesChannel(t *testing.T) {
	var buf bytes.Buffer
	n, _ := notify.New(&buf)
	ch := make(chan notify.Event, 3)
	ch <- notify.Event{Level: notify.LevelInfo, Port: 80, Protocol: "tcp", Message: "a"}
	ch <- notify.Event{Level: notify.LevelInfo, Port: 81, Protocol: "tcp", Message: "b"}
	ch <- notify.Event{Level: notify.LevelAlert, Port: 82, Protocol: "tcp", Message: "c"}
	close(ch)
	n.Watch(ch)
	got := buf.String()
	for _, want := range []string{"80/tcp", "81/tcp", "82/tcp"} {
		if !strings.Contains(got, want) {
			t.Errorf("expected %q in output, got: %q", want, got)
		}
	}
}
