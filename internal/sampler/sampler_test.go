package sampler

import (
	"context"
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

func stubScan(ports []scanner.Port) func(context.Context) ([]scanner.Port, error) {
	return func(_ context.Context) ([]scanner.Port, error) {
		return ports, nil
	}
}

func TestSamplerEmitsSample(t *testing.T) {
	ports := []scanner.Port{{Port: 80, Proto: "tcp"}}
	s := New(stubScan(ports), 20*time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	go s.Run(ctx)
	select {
	case sample, ok := <-s.Out():
		if !ok {
			t.Fatal("channel closed before sample")
		}
		if len(sample.Ports) != 1 || sample.Ports[0].Port != 80 {
			t.Fatalf("unexpected ports: %v", sample.Ports)
		}
	case <-ctx.Done():
		t.Fatal("timed out waiting for sample")
	}
}

func TestSamplerClosesChannelOnCancel(t *testing.T) {
	s := New(stubScan(nil), 10*time.Millisecond)
	ctx, cancel := context.WithCancel(context.Background())
	go s.Run(ctx)
	cancel()
	timer := time.NewTimer(200 * time.Millisecond)
	defer timer.Stop()
	for {
		select {
		case _, ok := <-s.Out():
			if !ok {
				return
			}
		case <-timer.C:
			t.Fatal("channel not closed after cancel")
		}
	}
}

func TestPipelineReturnsChannel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ch := NewPipeline(ctx, stubScan(nil), 50*time.Millisecond)
	if ch == nil {
		t.Fatal("expected non-nil channel")
	}
}
