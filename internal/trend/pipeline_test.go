package trend

import (
	"context"
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

func makePorts(n int) []scanner.Port {
	ports := make([]scanner.Port, n)
	for i := range ports {
		ports[i] = scanner.Port{Port: 8000 + i, Proto: "tcp"}
	}
	return ports
}

func TestPipelineRecordsTrend(t *testing.T) {
	tr := New(5)
	in := make(chan []scanner.Port, 4)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	out := NewPipeline(ctx, tr, in)

	in <- makePorts(3)
	in <- makePorts(6)

	var evs []Event
	timer := time.NewTimer(500 * time.Millisecond)
	defer timer.Stop()
	for len(evs) < 2 {
		select {
		case ev := <-out:
			evs = append(evs, ev)
		case <-timer.C:
			t.Fatal("timed out waiting for events")
		}
	}

	if evs[1].Dir != Rising {
		t.Fatalf("expected Rising, got %s", evs[1].Dir)
	}
	if evs[1].OpenCount != 6 {
		t.Fatalf("expected OpenCount 6, got %d", evs[1].OpenCount)
	}
}

func TestPipelineStopsOnContextCancel(t *testing.T) {
	tr := New(5)
	in := make(chan []scanner.Port)
	ctx, cancel := context.WithCancel(context.Background())
	out := NewPipeline(ctx, tr, in)
	cancel()
	select {
	case _, ok := <-out:
		if ok {
			t.Fatal("expected channel to be closed")
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timed out waiting for pipeline to stop")
	}
}

func TestPipelineStopsOnClosedInput(t *testing.T) {
	tr := New(5)
	in := make(chan []scanner.Port)
	ctx := context.Background()
	out := NewPipeline(ctx, tr, in)
	close(in)
	select {
	case _, ok := <-out:
		if ok {
			t.Fatal("expected channel to be closed")
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timed out waiting for pipeline to stop")
	}
}
