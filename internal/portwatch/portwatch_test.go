package portwatch_test

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/user/portwatch/internal/portwatch"
	"github.com/user/portwatch/internal/scanner"
)

func startListener(t *testing.T) (net.Listener, uint16) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	port := uint16(ln.Addr().(*net.TCPAddr).Port)
	return ln, port
}

func TestWatcherDetectsOpenedPort(t *testing.T) {
	ln, port := startListener(t)
	defer ln.Close()

	s := scanner.New(scanner.Options{
		Ports:    []uint16{port},
		Protocol: "tcp",
		Host:     "127.0.0.1",
	})
	w := portwatch.New(s, 50*time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	events := w.Watch(ctx)
	for e := range events {
		if e.State == "opened" && e.Port.Port == port {
			return
		}
	}
	t.Fatal("expected opened event for test port")
}

func TestWatcherDetectsClosedPort(t *testing.T) {
	ln, port := startListener(t)

	s := scanner.New(scanner.Options{
		Ports:    []uint16{port},
		Protocol: "tcp",
		Host:     "127.0.0.1",
	})
	w := portwatch.New(s, 50*time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	events := w.Watch(ctx)
	// Wait for the opened event first, then close the listener.
	for e := range events {
		if e.State == "opened" && e.Port.Port == port {
			ln.Close()
		}
		if e.State == "closed" && e.Port.Port == port {
			return
		}
	}
	t.Fatal("expected closed event for test port")
}

func TestPipelineForwardsEvents(t *testing.T) {
	ln, port := startListener(t)
	defer ln.Close()

	s := scanner.New(scanner.Options{
		Ports:    []uint16{port},
		Protocol: "tcp",
		Host:     "127.0.0.1",
	})
	cfg := portwatch.PipelineConfig{Interval: 50 * time.Millisecond}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	ch := portwatch.NewPipeline(ctx, s, cfg)
	for e := range ch {
		if e.State == "opened" && e.Port.Port == port {
			return
		}
	}
	t.Fatal("pipeline did not forward opened event")
}
