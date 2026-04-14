package scanner

import (
	"net"
	"testing"
	"time"
)

// startTestListener opens a TCP listener on a random port and returns it.
func startTestListener(t *testing.T) (net.Listener, int) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start test listener: %v", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	return ln, port
}

func TestScanDetectsOpenPort(t *testing.T) {
	ln, port := startTestListener(t)
	defer ln.Close()

	s := &Scanner{
		Host:    "127.0.0.1",
		MinPort: port,
		MaxPort: port,
		Timeout: 200 * time.Millisecond,
	}

	ports, err := s.Scan()
	if err != nil {
		t.Fatalf("Scan returned error: %v", err)
	}
	if len(ports) != 1 {
		t.Fatalf("expected 1 open port, got %d", len(ports))
	}
	if ports[0].Number != port {
		t.Errorf("expected port %d, got %d", port, ports[0].Number)
	}
	if ports[0].Protocol != "tcp" {
		t.Errorf("expected protocol tcp, got %s", ports[0].Protocol)
	}
}

func TestScanClosedPort(t *testing.T) {
	ln, port := startTestListener(t)
	ln.Close() // close immediately so the port is not open

	s := &Scanner{
		Host:    "127.0.0.1",
		MinPort: port,
		MaxPort: port,
		Timeout: 200 * time.Millisecond,
	}

	ports, err := s.Scan()
	if err != nil {
		t.Fatalf("Scan returned error: %v", err)
	}
	if len(ports) != 0 {
		t.Errorf("expected 0 open ports, got %d", len(ports))
	}
}

func TestPortString(t *testing.T) {
	p := Port{Number: 8080, Protocol: "tcp", Address: "127.0.0.1"}
	expected := "127.0.0.1:8080/tcp"
	if p.String() != expected {
		t.Errorf("expected %q, got %q", expected, p.String())
	}
}

func TestNewDefaults(t *testing.T) {
	s := New("localhost")
	if s.Host != "localhost" {
		t.Errorf("unexpected host: %s", s.Host)
	}
	if s.MinPort != 1 || s.MaxPort != 65535 {
		t.Errorf("unexpected port range: %d-%d", s.MinPort, s.MaxPort)
	}
	if s.Timeout != 500*time.Millisecond {
		t.Errorf("unexpected timeout: %v", s.Timeout)
	}
}
