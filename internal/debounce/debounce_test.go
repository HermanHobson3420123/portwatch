package debounce_test

import (
	"testing"
	"time"

	"portwatch/internal/debounce"
	"portwatch/internal/scanner"
)

func makePort(proto string, number int) scanner.Port {
	return scanner.Port{Proto: proto, Number: number}
}

func TestEventEmittedAfterDelay(t *testing.T) {
	d := debounce.New(50 * time.Millisecond)
	d.Push(debounce.Event{Port: makePort("tcp", 8080), Opened: true})

	select {
	case e := <-d.C():
		if !e.Opened || e.Port.Number != 8080 {
			t.Fatalf("unexpected event: %+v", e)
		}
	case <-time.After(300 * time.Millisecond):
		t.Fatal("expected event was not emitted")
	}
}

func TestOppositeEventCancelsPending(t *testing.T) {
	d := debounce.New(80 * time.Millisecond)
	d.Push(debounce.Event{Port: makePort("tcp", 9090), Opened: true})
	// Immediately close the same port — should cancel the open event.
	d.Push(debounce.Event{Port: makePort("tcp", 9090), Opened: false})

	select {
	case e := <-d.C():
		t.Fatalf("expected no event but got: %+v", e)
	case <-time.After(200 * time.Millisecond):
		// correct — nothing emitted
	}
}

func TestSameDirectionResetsTimer(t *testing.T) {
	d := debounce.New(80 * time.Millisecond)
	d.Push(debounce.Event{Port: makePort("tcp", 7070), Opened: true})
	time.Sleep(50 * time.Millisecond)
	// Push again before timer fires — should reset.
	d.Push(debounce.Event{Port: makePort("tcp", 7070), Opened: true})

	// Should not have fired yet (timer was reset).
	select {
	case <-d.C():
		t.Fatal("event fired too early")
	case <-time.After(40 * time.Millisecond):
	}

	// Now it should fire.
	select {
	case e := <-d.C():
		if e.Port.Number != 7070 {
			t.Fatalf("wrong port: %+v", e)
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatal("event not received after reset")
	}
}

func TestDifferentPortsAreIndependent(t *testing.T) {
	d := debounce.New(50 * time.Millisecond)
	d.Push(debounce.Event{Port: makePort("tcp", 1111), Opened: true})
	d.Push(debounce.Event{Port: makePort("tcp", 2222), Opened: true})

	received := map[int]bool{}
	timeout := time.After(300 * time.Millisecond)
	for len(received) < 2 {
		select {
		case e := <-d.C():
			received[e.Port.Number] = true
		case <-timeout:
			t.Fatalf("only received %d events", len(received))
		}
	}
}
