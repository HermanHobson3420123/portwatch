package portburst_test

import (
	"context"
	"testing"
	"time"

	"portwatch/internal/monitor"
	"portwatch/internal/portburst"
	"portwatch/internal/scanner"
)

func p(proto string, number int) scanner.Port {
	return scanner.Port{Proto: proto, Number: number}
}

func TestNoBurstBelowThreshold(t *testing.T) {
	d := portburst.New(time.Second, 3)
	port := p("tcp", 8080)
	for i := 0; i < 2; i++ {
		if ev := d.Record(port); ev != nil {
			t.Fatalf("unexpected burst event on event %d", i+1)
		}
	}
}

func TestBurstAtThreshold(t *testing.T) {
	d := portburst.New(time.Second, 3)
	port := p("tcp", 9090)
	var got *portburst.Event
	for i := 0; i < 3; i++ {
		got = d.Record(port)
	}
	if got == nil {
		t.Fatal("expected burst event, got nil")
	}
	if got.Count < 3 {
		t.Errorf("expected count >= 3, got %d", got.Count)
	}
	if got.Port != port {
		t.Errorf("expected port %v, got %v", port, got.Port)
	}
}

func TestResetClearsBurst(t *testing.T) {
	d := portburst.New(time.Second, 2)
	port := p("udp", 5353)
	d.Record(port)
	d.Record(port)
	d.Reset(port)
	if ev := d.Record(port); ev != nil {
		t.Fatal("expected no burst after reset")
	}
}

func TestDifferentPortsAreIndependent(t *testing.T) {
	d := portburst.New(time.Second, 3)
	a := p("tcp", 80)
	b := p("tcp", 443)
	d.Record(a)
	d.Record(a)
	if ev := d.Record(b); ev != nil {
		t.Fatal("port b should not trigger burst from port a events")
	}
}

func TestPipelineEmitsBurstEvent(t *testing.T) {
	d := portburst.New(time.Second, 2)
	in := make(chan monitor.Event, 4)
	port := p("tcp", 7777)

	for i := 0; i < 3; i++ {
		in <- monitor.Event{Port: port, Kind: monitor.Opened}
	}
	close(in)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	out := portburst.NewPipeline(ctx, d, in)
	var count int
	for range out {
		count++
	}
	if count == 0 {
		t.Fatal("expected at least one burst event from pipeline")
	}
}

func TestPipelineStopsOnContextCancel(t *testing.T) {
	d := portburst.New(time.Second, 2)
	in := make(chan monitor.Event)
	ctx, cancel := context.WithCancel(context.Background())
	out := portburst.NewPipeline(ctx, d, in)
	cancel()
	select {
	case _, ok := <-out:
		if ok {
			t.Fatal("expected channel to be closed after cancel")
		}
	case <-time.After(time.Second):
		t.Fatal("pipeline did not stop after context cancel")
	}
}
