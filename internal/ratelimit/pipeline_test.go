package ratelimit

import (
	"context"
	"testing"
	"time"

	"portwatch/internal/monitor"
	"portwatch/internal/scanner"
)

func makeEvent(port int) monitor.Event {
	return monitor.Event{
		Port:   scanner.Port{Port: port, Protocol: "tcp"},
		Opened: true,
	}
}

func TestPipelineForwardsFirstEvent(t *testing.T) {
	l := New(500 * time.Millisecond)
	in := make(chan monitor.Event, 4)
	p := NewPipeline(l, in)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go p.Run(ctx)

	in <- makeEvent(8080)
	select {
	case ev := <-p.Out():
		if ev.Port.Port != 8080 {
			t.Fatalf("expected port 8080, got %d", ev.Port.Port)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for event")
	}
}

func TestPipelineSuppressesDuplicateWithinCooldown(t *testing.T) {
	l := New(10 * time.Second)
	in := make(chan monitor.Event, 4)
	p := NewPipeline(l, in)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go p.Run(ctx)

	in <- makeEvent(9090)
	time.Sleep(50 * time.Millisecond)
	in <- makeEvent(9090) // should be suppressed

	count := 0
	timeout := time.After(300 * time.Millisecond)
loop:
	for {
		select {
		case <-p.Out():
			count++
		case <-timeout:
			break loop
		}
	}
	if count != 1 {
		t.Fatalf("expected 1 event, got %d", count)
	}
}

func TestPipelineStopsOnContextCancel(t *testing.T) {
	l := New(time.Second)
	in := make(chan monitor.Event)
	p := NewPipeline(l, in)
	ctx, cancel := context.WithCancel(context.Background())
	go p.Run(ctx)
	cancel()
	select {
	case _, ok := <-p.Out():
		if ok {
			t.Fatal("expected channel to be closed")
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for pipeline to stop")
	}
}
