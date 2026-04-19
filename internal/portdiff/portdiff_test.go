package portdiff_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"portwatch/internal/portdiff"
	"portwatch/internal/scanner"
)

func p(proto string, num int) scanner.Port {
	return scanner.Port{Protocol: proto, Number: num}
}

func TestCompareDetectsOpened(t *testing.T) {
	prev := []scanner.Port{p("tcp", 80)}
	next := []scanner.Port{p("tcp", 80), p("tcp", 443)}
	s := portdiff.Compare(prev, next)
	if len(s.Opened) != 1 || s.Opened[0].Number != 443 {
		t.Fatalf("expected 443 opened, got %v", s.Opened)
	}
	if len(s.Closed) != 0 {
		t.Fatalf("unexpected closed ports: %v", s.Closed)
	}
}

func TestCompareDetectsClosed(t *testing.T) {
	prev := []scanner.Port{p("tcp", 80), p("tcp", 8080)}
	next := []scanner.Port{p("tcp", 80)}
	s := portdiff.Compare(prev, next)
	if len(s.Closed) != 1 || s.Closed[0].Number != 8080 {
		t.Fatalf("expected 8080 closed, got %v", s.Closed)
	}
}

func TestSummaryEmpty(t *testing.T) {
	s := portdiff.Compare([]scanner.Port{p("tcp", 80)}, []scanner.Port{p("tcp", 80)})
	if !s.Empty() {
		t.Fatal("expected empty summary")
	}
}

func TestSummaryString(t *testing.T) {
	s := portdiff.Summary{
		Opened: []scanner.Port{p("tcp", 443)},
		Closed: []scanner.Port{p("tcp", 8080)},
	}
	out := s.String()
	if !strings.Contains(out, "+") || !strings.Contains(out, "-") {
		t.Fatalf("unexpected string output: %q", out)
	}
}

func TestPipelineEmitsDiff(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	in := make(chan []scanner.Port, 3)
	out := portdiff.NewPipeline(ctx, in)

	in <- []scanner.Port{p("tcp", 80)}
	in <- []scanner.Port{p("tcp", 80), p("tcp", 443)}
	close(in)

	select {
	case ev, ok := <-out:
		if !ok {
			t.Fatal("channel closed before event")
		}
		if len(ev.Summary.Opened) != 1 || ev.Summary.Opened[0].Number != 443 {
			t.Fatalf("unexpected event: %+v", ev)
		}
	case <-ctx.Done():
		t.Fatal("timeout waiting for diff event")
	}
}

func TestPipelineStopsOnContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	in := make(chan []scanner.Port)
	out := portdiff.NewPipeline(ctx, in)
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
