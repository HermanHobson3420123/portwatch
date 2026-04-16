package audit_test

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"portwatch/internal/audit"
)

func TestRecorderWritesEvents(t *testing.T) {
	dir := t.TempDir()
	l, _ := audit.New(filepath.Join(dir, "audit.json"))
	rec := audit.NewRecorder(l)

	ch := make(chan audit.PortEvent, 2)
	ch <- audit.PortEvent{Kind: audit.KindOpened, Port: makePort("tcp", 8080)}
	ch <- audit.PortEvent{Kind: audit.KindClosed, Port: makePort("tcp", 8080)}
	close(ch)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	rec.Watch(ctx, ch)

	evs := l.Events()
	if len(evs) != 2 {
		t.Fatalf("expected 2 events, got %d", len(evs))
	}
	if evs[0].Kind != audit.KindOpened {
		t.Errorf("first kind = %s", evs[0].Kind)
	}
	if evs[1].Kind != audit.KindClosed {
		t.Errorf("second kind = %s", evs[1].Kind)
	}
}

func TestRecorderStopsOnContextCancel(t *testing.T) {
	dir := t.TempDir()
	l, _ := audit.New(filepath.Join(dir, "audit.json"))
	rec := audit.NewRecorder(l)

	ch := make(chan audit.PortEvent)
	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan struct{})
	go func() {
		rec.Watch(ctx, ch)
		close(done)
	}()

	cancel()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Error("Watch did not stop after context cancel")
	}
}
