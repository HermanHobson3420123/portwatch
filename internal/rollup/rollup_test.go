package rollup_test

import (
	"context"
	"testing"
	"time"

	"portwatch/internal/rollup"
	"portwatch/internal/scanner"
)

func makePort(proto string, number int) scanner.Port {
	return scanner.Port{Proto: proto, Number: number}
}

func TestFlushEmitsBatchedEvents(t *testing.T) {
	r := rollup.New(30 * time.Millisecond)
	r.AddOpened(makePort("tcp", 8080))
	r.AddOpened(makePort("tcp", 9090))
	r.AddClosed(makePort("udp", 53))

	select {
	case ev := <-r.Events():
		if len(ev.Opened) != 2 {
			t.Fatalf("expected 2 opened, got %d", len(ev.Opened))
		}
		if len(ev.Closed) != 1 {
			t.Fatalf("expected 1 closed, got %d", len(ev.Closed))
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatal("timed out waiting for rollup event")
	}
}

func TestOppositeEventsCancel(t *testing.T) {
	r := rollup.New(30 * time.Millisecond)
	r.AddOpened(makePort("tcp", 8080))
	r.AddClosed(makePort("tcp", 8080)) // cancels the open

	select {
	case ev := <-r.Events():
		if len(ev.Opened) != 0 {
			t.Fatalf("expected 0 opened after cancel, got %d", len(ev.Opened))
		}
		if len(ev.Closed) != 1 {
			t.Fatalf("expected 1 closed, got %d", len(ev.Closed))
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatal("timed out")
	}
}

func TestTimerResetsOnNewEvent(t *testing.T) {
	r := rollup.New(40 * time.Millisecond)
	start := time.Now()
	r.AddOpened(makePort("tcp", 1111))
	time.Sleep(20 * time.Millisecond)
	r.AddOpened(makePort("tcp", 2222)) // resets the timer

	select {
	case ev := <-r.Events():
		elapsed := time.Since(start)
		if elapsed < 50*time.Millisecond {
			t.Fatalf("event emitted too early: %v", elapsed)
		}
		if len(ev.Opened) != 2 {
			t.Fatalf("expected 2 opened, got %d", len(ev.Opened))
		}
	case <-time.After(300 * time.Millisecond):
		t.Fatal("timed out")
	}
}

func TestWatchConsumesEvents(t *testing.T) {
	r := rollup.New(20 * time.Millisecond)
	r.AddOpened(makePort("tcp", 443))

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	received := make(chan rollup.Event, 1)
	go r.Watch(ctx, func(ev rollup.Event) {
		received <- ev
	})

	select {
	case ev := <-received:
		if len(ev.Opened) != 1 {
			t.Fatalf("expected 1 opened, got %d", len(ev.Opened))
		}
	case <-ctx.Done():
		t.Fatal("timed out waiting for Watch callback")
	}
}
