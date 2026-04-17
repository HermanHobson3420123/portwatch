package sampler

import (
	"context"
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

func TestPipelineForwardsSamples(t *testing.T) {
	ports := []scanner.Port{{Port: 443, Proto: "tcp"}}
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()
	ch := NewPipeline(ctx, stubScan(ports), 20*time.Millisecond)
	select {
	case s, ok := <-ch:
		if !ok {
			t.Fatal("channel closed unexpectedly")
		}
		if len(s.Ports) == 0 {
			t.Fatal("expected ports in sample")
		}
	case <-ctx.Done():
		t.Fatal("timed out waiting for pipeline sample")
	}
}

func TestPipelineStopsOnContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	ch := NewPipeline(ctx, stubScan(nil), 10*time.Millisecond)
	cancel()
	timer := time.NewTimer(200 * time.Millisecond)
	defer timer.Stop()
	for {
		select {
		case _, ok := <-ch:
			if !ok {
				return
			}
		case <-timer.C:
			t.Fatal("pipeline did not stop after context cancel")
		}
	}
}
