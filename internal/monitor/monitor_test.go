package monitor_test

import (
	"net"
	"testing"
	"time"

	"github.com/user/portwatch/internal/monitor"
	"github.com/user/portwatch/internal/scanner"
)

func startListener(t *testing.T) (net.Listener, int) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start listener: %v", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	return ln, port
}

func TestMonitorDetectsOpenedPort(t *testing.T) {
	s := scanner.New("127.0.0.1", []int{})

	// Start monitor with no ports open
	m := monitor.New(s, 50*time.Millisecond)
	m.Start()
	defer m.Stop()

	// Give baseline scan time to complete
	time.Sleep(80 * time.Millisecond)

	// Open a port after baseline
	ln, port := startListener(t)
	defer ln.Close()

	s2 := scanner.New("127.0.0.1", []int{port})
	m2 := monitor.New(s2, 50*time.Millisecond)
	m2.Start()
	defer m2.Stop()

	time.Sleep(80 * time.Millisecond)

	select {
	case ch := <-m2.Changes:
		if !ch.Opened {
			t.Errorf("expected Opened=true, got false")
		}
	case <-time.After(500 * time.Millisecond):
		t.Error("timed out waiting for opened-port change")
	}
}

func TestMonitorDetectsClosedPort(t *testing.T) {
	ln, port := startListener(t)

	s := scanner.New("127.0.0.1", []int{port})
	m := monitor.New(s, 50*time.Millisecond)
	m.Start()
	defer m.Stop()

	// Let baseline establish with port open
	time.Sleep(80 * time.Millisecond)

	// Close the port
	ln.Close()

	select {
	case ch := <-m.Changes:
		if ch.Opened {
			t.Errorf("expected Opened=false (closed), got true")
		}
	case <-time.After(500 * time.Millisecond):
		t.Error("timed out waiting for closed-port change")
	}
}
