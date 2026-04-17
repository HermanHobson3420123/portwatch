package trend

import (
	"sync"
	"time"
)

// Direction indicates whether port activity is increasing or decreasing.
type Direction int

const (
	Stable Direction = iota
	Rising
	Falling
)

func (d Direction) String() string {
	switch d {
	case Rising:
		return "rising"
	case Falling:
		return "falling"
	default:
		return "stable"
	}
}

// Sample holds an open-port count at a point in time.
type Sample struct {
	At    time.Time
	Count int
}

// Tracker accumulates port-count samples and derives a trend direction.
type Tracker struct {
	mu      sync.Mutex
	window  int
	samples []Sample
}

// New returns a Tracker that keeps up to window samples.
func New(window int) *Tracker {
	if window < 2 {
		window = 2
	}
	return &Tracker{window: window}
}

// Record adds a new sample.
func (t *Tracker) Record(count int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.samples = append(t.samples, Sample{At: time.Now(), Count: count})
	if len(t.samples) > t.window {
		t.samples = t.samples[len(t.samples)-t.window:]
	}
}

// Direction returns the current trend based on the first and last sample.
func (t *Tracker) Direction() Direction {
	t.mu.Lock()
	defer t.mu.Unlock()
	if len(t.samples) < 2 {
		return Stable
	}
	first := t.samples[0].Count
	last := t.samples[len(t.samples)-1].Count
	switch {
	case last > first:
		return Rising
	case last < first:
		return Falling
	default:
		return Stable
	}
}

// Samples returns a copy of the current sample slice.
func (t *Tracker) Samples() []Sample {
	t.mu.Lock()
	defer t.mu.Unlock()
	out := make([]Sample, len(t.samples))
	copy(out, t.samples)
	return out
}
