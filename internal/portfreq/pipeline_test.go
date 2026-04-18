package portfreq

import (
	"context"
	"testing"

	"portwatch/internal/scanner"
)

func TestPipelineRecordsAndForwards(t *testing.T) {
	tr := New()
	in := make(chan []scanner.Port, 1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	out := NewPipeline(ctx, tr, in)

	in <- []scanner.Port{p("tcp", "0.0.0.0:80"), p("tcp", "0.0.0.0:443")}
	close(in)

	ports := <-out
	if len(ports) != 2 {
		t.Fatalf("expected 2 ports forwarded, got %d", len(ports))
	}
	if len(tr.Top(10)) != 2 {
		t.Fatalf("expected 2 entries in tracker, got %d", len(tr.Top(10)))
	}
}

func TestPipelineStopsOnContextCancel(t *testing.T) {
	tr := New()
	in := make(chan []scanner.Port)
	ctx, cancel := context.WithCancel(context.Background())
	out := NewPipeline(ctx, tr, in)
	cancel()
	_, ok := <-out
	if ok {
		t.Fatal("expected pipeline to close after context cancel")
	}
}

func TestPipelineStopsOnClosedInput(t *testing.T) {
	tr := New()
	in := make(chan []scanner.Port)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	out := NewPipeline(ctx, tr, in)
	close(in)
	_, ok := <-out
	if ok {
		t.Fatal("expected pipeline to close when input closes")
	}
}
