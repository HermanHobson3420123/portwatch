// Package debounce delays forwarding of port events to reduce noise from
// transient port fluctuations (e.g. short-lived ephemeral connections).
package debounce

import (
	"sync"
	"time"

	"portwatch/internal/scanner"
)

// Event wraps a port change with its direction.
type Event struct {
	Port    scanner.Port
	Opened  bool
}

// Debouncer holds pending events and emits them only after a quiet period.
type Debouncer struct {
	delay   time.Duration
	mu      sync.Mutex
	pending map[string]*entry
	out     chan Event
}

type entry struct {
	event Event
	timer *time.Timer
}

// New creates a Debouncer that waits delay before forwarding each event.
func New(delay time.Duration) *Debouncer {
	return &Debouncer{
		delay:   delay,
		pending: make(map[string]*entry),
		out:     make(chan Event, 64),
	}
}

// C returns the output channel of debounced events.
func (d *Debouncer) C() <-chan Event {
	return d.out
}

// Push schedules an event for emission after the debounce delay.
// If a conflicting event for the same port arrives before the delay
// expires, the pending event is cancelled.
func (d *Debouncer) Push(e Event) {
	key := portKey(e.Port)

	d.mu.Lock()
	defer d.mu.Unlock()

	if existing, ok := d.pending[key]; ok {
		// Opposite direction cancels the pending event entirely.
		if existing.event.Opened != e.Opened {
			existing.timer.Stop()
			delete(d.pending, key)
			return
		}
		// Same direction: reset the timer.
		existing.timer.Reset(d.delay)
		return
	}

	en := &entry{event: e}
	en.timer = time.AfterFunc(d.delay, func() {
		d.mu.Lock()
		delete(d.pending, key)
		d.mu.Unlock()
		d.out <- e
	})
	d.pending[key] = en
}

func portKey(p scanner.Port) string {
	return p.Proto + ":" + p.String()
}
