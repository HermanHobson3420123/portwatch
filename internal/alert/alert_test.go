package alert_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/monitor"
	"github.com/user/portwatch/internal/scanner"
)

func TestNotifyOpened(t *testing.T) {
	var buf bytes.Buffer
	a := alert.New(&buf, alert.LevelInfo)

	c := monitor.Change{
		Port:   scanner.Port{Address: "127.0.0.1:8080"},
		Opened: true,
	}
	a.Notify(c)

	out := buf.String()
	if !strings.Contains(out, "OPENED") {
		t.Errorf("expected OPENED in output, got: %s", out)
	}
	if !strings.Contains(out, "127.0.0.1:8080") {
		t.Errorf("expected address in output, got: %s", out)
	}
	if !strings.Contains(out, "[+]") {
		t.Errorf("expected [+] symbol in output, got: %s", out)
	}
}

func TestNotifyClosed(t *testing.T) {
	var buf bytes.Buffer
	a := alert.New(&buf, alert.LevelWarn)

	c := monitor.Change{
		Port:   scanner.Port{Address: "127.0.0.1:9090"},
		Opened: false,
	}
	a.Notify(c)

	out := buf.String()
	if !strings.Contains(out, "CLOSED") {
		t.Errorf("expected CLOSED in output, got: %s", out)
	}
	if !strings.Contains(out, "[-]") {
		t.Errorf("expected [-] symbol in output, got: %s", out)
	}
}

func TestWatchConsumesChannel(t *testing.T) {
	var buf bytes.Buffer
	a := alert.New(&buf, alert.LevelInfo)

	ch := make(chan monitor.Change, 3)
	ch <- monitor.Change{Port: scanner.Port{Address: "127.0.0.1:1111"}, Opened: true}
	ch <- monitor.Change{Port: scanner.Port{Address: "127.0.0.1:2222"}, Opened: false}
	ch <- monitor.Change{Port: scanner.Port{Address: "127.0.0.1:3333"}, Opened: true}
	close(ch)

	a.Watch(ch)

	out := buf.String()
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 3 {
		t.Errorf("expected 3 lines of output, got %d: %s", len(lines), out)
	}
}
