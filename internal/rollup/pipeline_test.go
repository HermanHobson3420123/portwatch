package rollup_test

import (
	"context"
	"testing"
	"time"

	"portwatch/internal/monitor"
	"portwatch/internal/rollup"
	"portwatch/internal/scanner"
)

func TestPipelineForwardsDiff(t *testing.T) {
	r := rollup.New(20 * time.Millisecond)
	pl := rollup.NewPipeline(r)

	diffs := make(chan monitor.Diff, 1)
	diffs <- monitor.Diff{
		Opened: []scanner.Port{makePort("tcp", 8443)},
		Closed: []scanner.Port{makePort("tcp", 80)},
	}
	close(diffs)

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	received := make(chan rollup.Event, 1)
	go pl.Run(ctx, diffs, func(ev rollup.Event) {
		received <- ev
	})

	select {
	case ev := <-received:
		if len(ev.Opened) != 1 || ev.Opened[0].Number != 8443 {
			t.Fatalf("unexpected opened ports: %+v", ev.Opened)
		}
		if len(ev.Closed) != 1 || ev.Closed[0].Number != 80 {
			t.Fatalf("unexpected closed ports: %+v", ev.Closed)
		}
	case <-ctx.Done():
		t.Fatal("timed out waiting for pipeline event")
	}
}

func TestPipelineStopsOnContextCancel(t *testing.T) {
	r := rollup.New(50 * time.Millisecond)
	pl := rollup.NewPipeline(r)

	diffs := make(chan monitor.Diff)
	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan struct{})
	go func() {
		pl.Run(ctx, diffs, func(rollup.Event) {})
		close(done)
	}()

	cancel()
	select {
	case <-done:
	case <-time.After(200 * time.Millisecond):
		t.Fatal("pipeline did not stop after context cancel")
	}
}
