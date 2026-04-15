package debounce_test

import (
	"context"
	"testing"
	"time"

	"portwatch/internal/debounce"
	"portwatch/internal/scanner"
)

func TestPipelineForwardsEvent(t *testing.T) {
	in := make(chan debounce.PortChange, 4)
	pl := debounce.NewPipeline(40*time.Millisecond, in)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go pl.Run(ctx)

	in <- debounce.PortChange{Port: scanner.Port{Proto: "tcp", Number: 3000}, Opened: true}

	select {
	case e := <-pl.C():
		if e.Port.Number != 3000 || !e.Opened {
			t.Fatalf("unexpected event: %+v", e)
		}
	case <-time.After(300 * time.Millisecond):
		t.Fatal("pipeline did not emit event")
	}
}

func TestPipelineStopsOnContextCancel(t *testing.T) {
	in := make(chan debounce.PortChange)
	pl := debounce.NewPipeline(40*time.Millisecond, in)

	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan struct{})
	go func() {
		pl.Run(ctx)
		close(done)
	}()

	cancel()

	select {
	case <-done:
		// OK
	case <-time.After(200 * time.Millisecond):
		t.Fatal("pipeline did not stop after context cancellation")
	}
}

func TestPipelineStopsOnClosedChannel(t *testing.T) {
	in := make(chan debounce.PortChange)
	pl := debounce.NewPipeline(40*time.Millisecond, in)

	ctx := context.Background()
	done := make(chan struct{})
	go func() {
		pl.Run(ctx)
		close(done)
	}()

	close(in)

	select {
	case <-done:
		// OK
	case <-time.After(200 * time.Millisecond):
		t.Fatal("pipeline did not stop after input channel closed")
	}
}
