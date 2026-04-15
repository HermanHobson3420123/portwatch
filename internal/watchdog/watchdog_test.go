package watchdog_test

import (
	"context"
	"testing"
	"time"

	"portwatch/internal/scanner"
	"portwatch/internal/watchdog"
)

func newScanner(t *testing.T) *scanner.Scanner {
	t.Helper()
	s, err := scanner.New("tcp", []string{"127.0.0.1"}, 1*time.Second)
	if err != nil {
		t.Fatalf("scanner.New: %v", err)
	}
	return s
}

func TestBeatUpdatesStatus(t *testing.T) {
	s := newScanner(t)
	wd := watchdog.New(s, 50*time.Millisecond, 200*time.Millisecond)

	if wd.Status().ScanCount != 0 {
		t.Fatal("expected zero scan count initially")
	}

	wd.Beat()
	wd.Beat()

	st := wd.Status()
	if st.ScanCount != 2 {
		t.Fatalf("expected ScanCount=2, got %d", st.ScanCount)
	}
	if st.LastScan.IsZero() {
		t.Fatal("LastScan should not be zero after Beat")
	}
}

func TestWatchDetectsStaleScanner(t *testing.T) {
	s := newScanner(t)
	wd := watchdog.New(s, 20*time.Millisecond, 30*time.Millisecond)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go wd.Watch(ctx)

	// Send one beat then go silent.
	wd.Beat()

	// Wait long enough for the timeout to trigger.
	time.Sleep(120 * time.Millisecond)

	if wd.Status().Healthy {
		t.Fatal("expected watchdog to mark scanner as unhealthy after timeout")
	}
}

func TestWatchHealthyWhileBeating(t *testing.T) {
	s := newScanner(t)
	wd := watchdog.New(s, 20*time.Millisecond, 80*time.Millisecond)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go wd.Watch(ctx)

	// Keep sending beats faster than the timeout.
	for i := 0; i < 6; i++ {
		wd.Beat()
		time.Sleep(15 * time.Millisecond)
	}

	if !wd.Status().Healthy {
		t.Fatal("expected watchdog to remain healthy while receiving beats")
	}
}

func TestWatchExitsOnContextCancel(t *testing.T) {
	s := newScanner(t)
	wd := watchdog.New(s, 10*time.Millisecond, 50*time.Millisecond)

	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan struct{})
	go func() {
		wd.Watch(ctx)
		close(done)
	}()

	cancel()

	select {
	case <-done:
		// ok
	case <-time.After(200 * time.Millisecond):
		t.Fatal("Watch did not exit after context cancel")
	}
}
