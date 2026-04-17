// Package throttle limits how frequently scan cycles are allowed to run,
// preventing CPU spikes when the poll interval is set very low.
package throttle

import (
	"sync"
	"time"
)

// Throttle enforces a minimum duration between successive ticks.
type Throttle struct {
	mu       sync.Mutex
	minGap   time.Duration
	lastTick time.Time
	now      func() time.Time
}

// New returns a Throttle that enforces at least minGap between ticks.
// If minGap is zero or negative the throttle is effectively disabled.
func New(minGap time.Duration) *Throttle {
	return &Throttle{
		minGap: minGap,
		now:    time.Now,
	}
}

// Allow returns true if enough time has elapsed since the last allowed tick.
// If allowed, the internal timestamp is updated.
func (t *Throttle) Allow() bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := t.now()
	if t.minGap <= 0 || t.lastTick.IsZero() || now.Sub(t.lastTick) >= t.minGap {
		t.lastTick = now
		return true
	}
	return false
}

// Wait blocks until the throttle would allow a tick, then marks it.
func (t *Throttle) Wait() {
	for {
		t.mu.Lock()
		now := t.now()
		if t.minGap <= 0 || t.lastTick.IsZero() || now.Sub(t.lastTick) >= t.minGap {
			t.lastTick = now
			t.mu.Unlock()
			return
		}
		remaining := t.minGap - now.Sub(t.lastTick)
		t.mu.Unlock()
		time.Sleep(remaining)
	}
}

// Reset clears the last-tick timestamp so the next call to Allow or Wait
// succeeds immediately.
func (t *Throttle) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.lastTick = time.Time{}
}

// Remaining returns how much time must still elapse before the next tick is
// allowed. Returns zero if a tick is allowed immediately.
func (t *Throttle) Remaining() time.Duration {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.minGap <= 0 || t.lastTick.IsZero() {
		return 0
	}
	elapsed := t.now().Sub(t.lastTick)
	if elapsed >= t.minGap {
		return 0
	}
	return t.minGap - elapsed
}
