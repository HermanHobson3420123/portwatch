package portstate

import (
	"context"
	"testing"
	"time"

	"portwatch/internal/monitor"
	"portwatch/internal/scanner"
)

func makeEvent(typ monitor.EventType, port int) monitor.Event {
	return monitor.Event{Type: typ, Port: scanner.Port{Proto: "tcp", Number: port}}
}

func TestPipelineTracksOpenedPort(t *testing.T) {
	tr := New()
	in := make(chan monitor.Event, 1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	out := NewPipeline(ctx, tr, in)

	in <- makeEvent(monitor.Opened, 8080)
	ev := <-out
	if ev.Type != monitor.Opened {
		t.Fatalf("expected Opened, got %v", ev.Type)
	}
	if tr.Len() != 1 {
		t.Fatalf("tracker should have 1 entry, got %d", tr.Len())
	}
}

func TestPipelineTracksClosedPort(t *testing.T) {
	tr := New()
	tr.Open(scanner.Port{Proto: "tcp", Number: 9090}, time.Now())
	in := make(chan monitor.Event, 1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	out := NewPipeline(ctx, tr, in)

	in <- makeEvent(monitor.Closed, 9090)
	<-out
	if tr.Len() != 0 {
		t.Fatalf("tracker should be empty after close, got %d", tr.Len())
	}
}

func TestPipelineStopsOnContextCancel(t *testing.T) {
	tr := New()
	in := make(chan monitor.Event)
	ctx, cancel := context.WithCancel(context.Background())
	out := NewPipeline(ctx, tr, in)
	cancel()
	select {
	case _, ok := <-out:
		if ok {
			t.Fatal("expected channel to be closed")
		}
	case <-time.After(time.Second):
		t.Fatal("pipeline did not stop after context cancel")
	}
}

func TestPipelineStopsOnClosedInput(t *testing.T) {
	tr := New()
	in := make(chan monitor.Event)
	ctx := context.Background()
	out := NewPipeline(ctx, tr, in)
	close(in)
	select {
	case _, ok := <-out:
		if ok {
			t.Fatal("expected channel to be closed")
		}
	case <-time.After(time.Second):
		t.Fatal("pipeline did not stop after input closed")
	}
}
