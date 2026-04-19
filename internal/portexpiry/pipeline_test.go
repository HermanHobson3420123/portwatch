package portexpiry

import (
	"context"
	"testing"
	"time"

	"portwatch/internal/scanner"
)

func TestPipelineEmitsExpiredPort(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	in := make(chan scanner.Port, 1)
	in <- scanner.Port{Proto: "tcp", Number: 9090}
	close(in)

	events := NewPipeline(ctx, in, 50*time.Millisecond, 30*time.Millisecond)

	select {
	case ev, ok := <-events:
		if !ok {
			t.Fatal("channel closed before expiry event")
		}
		if ev.Port.Number != 9090 {
			t.Errorf("expected port 9090, got %d", ev.Port.Number)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timed out waiting for expiry event")
	}
}

func TestPipelineStopsOnContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	in := make(chan scanner.Port)
	events := NewPipeline(ctx, in, time.Minute, time.Minute)
	cancel()
	select {
	case <-events:
	case <-time.After(200 * time.Millisecond):
		t.Fatal("pipeline did not stop after context cancel")
	}
}
