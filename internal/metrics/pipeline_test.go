package metrics

import (
	"context"
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

func makeEvent(opened bool, port uint16) Event {
	return Event{Opened: opened, Port: scanner.Port{Port: port, Proto: "tcp"}}
}

func TestPipelineRecordsAndForwards(t *testing.T) {
	c := New()
	in := make(chan Event, 4)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	out := NewPipeline(ctx, c, in)

	in <- makeEvent(true, 80)
	in <- makeEvent(false, 443)

	for i := 0; i < 2; i++ {
		select {
		case <-out:
		case <-time.After(time.Second):
			t.Fatal("timeout waiting for event")
		}
	}

	s := c.Snapshot()
	if s.OpenedTotal != 1 {
		t.Fatalf("expected 1 opened, got %d", s.OpenedTotal)
	}
	if s.ClosedTotal != 1 {
		t.Fatalf("expected 1 closed, got %d", s.ClosedTotal)
	}
}

func TestPipelineStopsOnContextCancel(t *testing.T) {
	c := New()
	in := make(chan Event)
	ctx, cancel := context.WithCancel(context.Background())
	out := NewPipeline(ctx, c, in)
	cancel()
	select {
	case _, ok := <-out:
		if ok {
			t.Fatal("expected closed channel")
		}
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for pipeline to stop")
	}
}

func TestPipelineStopsOnClosedInput(t *testing.T) {
	c := New()
	in := make(chan Event)
	ctx := context.Background()
	out := NewPipeline(ctx, c, in)
	close(in)
	select {
	case _, ok := <-out:
		if ok {
			t.Fatal("expected closed channel")
		}
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for pipeline to stop")
	}
}
