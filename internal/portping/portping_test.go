package portping_test

import (
	"context"
	"net"
	"testing"
	"time"

	"portwatch/internal/portping"
	"portwatch/internal/scanner"
)

func startTCPListener(t *testing.T) (uint16, func()) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	addr := ln.Addr().(*net.TCPAddr)
	return uint16(addr.Port), func() { _ = ln.Close() }
}

func TestPingAlive(t *testing.T) {
	port, stop := startTCPListener(t)
	defer stop()

	p := portping.New(time.Second)
	res := p.Ping(context.Background(), scanner.Port{Number: port, Protocol: "tcp"})
	if !res.Alive {
		t.Fatalf("expected alive, got err: %v", res.Err)
	}
	if res.Latency <= 0 {
		t.Fatal("expected positive latency")
	}
}

func TestPingClosed(t *testing.T) {
	p := portping.New(200 * time.Millisecond)
	res := p.Ping(context.Background(), scanner.Port{Number: 1, Protocol: "tcp"})
	if res.Alive {
		t.Fatal("expected not alive for port 1")
	}
}

func TestPingAllReturnsAllResults(t *testing.T) {
	port, stop := startTCPListener(t)
	defer stop()

	ports := []scanner.Port{
		{Number: port, Protocol: "tcp"},
		{Number: 1, Protocol: "tcp"},
	}
	p := portping.New(200 * time.Millisecond)
	results := p.PingAll(context.Background(), ports)
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestPipelineForwardsResults(t *testing.T) {
	port, stop := startTCPListener(t)
	defer stop()

	in := make(chan []scanner.Port, 1)
	in <- []scanner.Port{{Number: port, Protocol: "tcp"}}
	close(in)

	p := portping.New(time.Second)
	ctx := context.Background()
	out := portping.NewPipeline(ctx, p, in)

	results, ok := <-out
	if !ok {
		t.Fatal("expected results from pipeline")
	}
	if len(results) != 1 || !results[0].Alive {
		t.Fatalf("unexpected results: %+v", results)
	}
}

func TestPipelineStopsOnContextCancel(t *testing.T) {
	in := make(chan []scanner.Port)
	p := portping.New(time.Second)
	ctx, cancel := context.WithCancel(context.Background())
	out := portping.NewPipeline(ctx, p, in)
	cancel()
	_, ok := <-out
	if ok {
		t.Fatal("expected pipeline to stop")
	}
}
