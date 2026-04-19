package portschedule

import (
	"context"
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

func ports(nums ...int) []scanner.Port {
	var out []scanner.Port
	for _, n := range nums {
		out = append(out, scanner.Port{Number: n, Proto: "tcp"})
	}
	return out
}

func TestGateNoActiveWindowPassesAll(t *testing.T) {
	s := New()
	g := NewGate(s)
	result := g.Filter(ports(80, 443, 8080))
	if len(result) != 3 {
		t.Errorf("expected 3, got %d", len(result))
	}
}

func TestGateFiltersToWindow(t *testing.T) {
	s := New()
	now := time.Now()
	s.Add(Entry{Label: "web", Start: now.Add(-time.Minute), End: now.Add(time.Minute), PortLow: 80, PortHigh: 80, Protocol: "tcp"})
	g := NewGate(s)
	result := g.Filter(ports(80, 443))
	if len(result) != 1 || result[0].Number != 80 {
		t.Errorf("expected only port 80, got %+v", result)
	}
}

func TestPipelineForwardsPorts(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s := New()
	g := NewGate(s)
	in := make(chan []scanner.Port, 1)
	out := NewPipeline(ctx, g, in)

	in <- ports(22, 80)
	got := <-out
	if len(got) != 2 {
		t.Errorf("expected 2 ports, got %d", len(got))
	}
}

func TestPipelineStopsOnContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	s := New()
	g := NewGate(s)
	in := make(chan []scanner.Port)
	out := NewPipeline(ctx, g, in)
	cancel()
	_, ok := <-out
	if ok {
		t.Error("expected pipeline to close after context cancel")
	}
}

func TestPipelineStopsOnClosedInput(t *testing.T) {
	ctx := context.Background()
	s := New()
	g := NewGate(s)
	in := make(chan []scanner.Port)
	out := NewPipeline(ctx, g, in)
	close(in)
	_, ok := <-out
	if ok {
		t.Error("expected pipeline to close when input closes")
	}
}
