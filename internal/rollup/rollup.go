// Package rollup batches rapid-fire port change events into a single
// summary notification after a quiet period, reducing alert noise during
// large network topology changes.
package rollup

import (
	"context"
	"sync"
	"time"

	"portwatch/internal/scanner"
)

// Event holds a batched summary of port changes.
type Event struct {
	Opened []scanner.Port
	Closed []scanner.Port
	At     time.Time
}

// Rollup collects individual open/close events and emits a combined Event
// once no new events arrive within the quiet window.
type Rollup struct {
	window time.Duration
	mu     sync.Mutex
	opened map[string]scanner.Port
	closed map[string]scanner.Port
	timer  *time.Timer
	out    chan Event
}

// New creates a Rollup that waits window duration after the last event
// before emitting a batched summary.
func New(window time.Duration) *Rollup {
	return &Rollup{
		window: window,
		opened: make(map[string]scanner.Port),
		closed: make(map[string]scanner.Port),
		out:    make(chan Event, 8),
	}
}

// AddOpened records a newly opened port.
func (r *Rollup) AddOpened(p scanner.Port) {
	r.mu.Lock()
	defer r.mu.Unlock()
	key := portKey(p)
	delete(r.closed, key)
	r.opened[key] = p
	r.resetTimer()
}

// AddClosed records a newly closed port.
func (r *Rollup) AddClosed(p scanner.Port) {
	r.mu.Lock()
	defer r.mu.Unlock()
	key := portKey(p)
	delete(r.opened, key)
	r.closed[key] = p
	r.resetTimer()
}

// Events returns the channel on which batched Events are delivered.
func (r *Rollup) Events() <-chan Event {
	return r.out
}

// Watch drains the output channel until ctx is cancelled.
func (r *Rollup) Watch(ctx context.Context, fn func(Event)) {
	for {
		select {
		case <-ctx.Done():
			return
		case ev, ok := <-r.out:
			if !ok {
				return
			}
			fn(ev)
		}
	}
}

// resetTimer (re)starts the quiet-period timer; must be called with r.mu held.
func (r *Rollup) resetTimer() {
	if r.timer != nil {
		r.timer.Stop()
	}
	r.timer = time.AfterFunc(r.window, r.flush)
}

func (r *Rollup) flush() {
	r.mu.Lock()
	defer r.mu.Unlock()
	if len(r.opened) == 0 && len(r.closed) == 0 {
		return
	}
	ev := Event{At: time.Now()}
	for _, p := range r.opened {
		ev.Opened = append(ev.Opened, p)
	}
	for _, p := range r.closed {
		ev.Closed = append(ev.Closed, p)
	}
	r.opened = make(map[string]scanner.Port)
	r.closed = make(map[string]scanner.Port)
	r.out <- ev
}

func portKey(p scanner.Port) string {
	return p.Proto + ":" + p.String()
}
