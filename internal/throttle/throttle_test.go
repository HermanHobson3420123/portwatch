package throttle

import (
	"testing"
	"time"
)

func fixedClock(base time.Time) func() time.Time {
	current := base
	return func() time.Time { return current }
}

func TestAllowFirstTickAlwaysPasses(t *testing.T) {
	th := New(100 * time.Millisecond)
	if !th.Allow() {
		t.Fatal("expected first tick to be allowed")
	}
}

func TestAllowSecondTickSuppressedWithinGap(t *testing.T) {
	base := time.Now()
	th := New(100 * time.Millisecond)
	th.now = func() time.Time { return base }

	th.Allow() // consume the first tick

	// Advance only 50 ms — still within the gap.
	th.now = func() time.Time { return base.Add(50 * time.Millisecond) }
	if th.Allow() {
		t.Fatal("expected second tick to be suppressed within gap")
	}
}

func TestAllowPassesAfterGapElapsed(t *testing.T) {
	base := time.Now()
	th := New(100 * time.Millisecond)
	th.now = func() time.Time { return base }

	th.Allow()

	th.now = func() time.Time { return base.Add(100 * time.Millisecond) }
	if !th.Allow() {
		t.Fatal("expected tick to pass after gap elapsed")
	}
}

func TestZeroGapAlwaysAllows(t *testing.T) {
	th := New(0)
	for i := 0; i < 5; i++ {
		if !th.Allow() {
			t.Fatalf("expected allow on iteration %d with zero gap", i)
		}
	}
}

func TestResetAllowsImmediateTick(t *testing.T) {
	base := time.Now()
	th := New(500 * time.Millisecond)
	th.now = func() time.Time { return base }

	th.Allow()

	// Still within gap — would normally be suppressed.
	th.now = func() time.Time { return base.Add(10 * time.Millisecond) }
	th.Reset()

	if !th.Allow() {
		t.Fatal("expected Allow to pass immediately after Reset")
	}
}

func TestWaitReturnsAfterGap(t *testing.T) {
	th := New(20 * time.Millisecond)
	th.Allow() // set lastTick to now

	start := time.Now()
	th.Wait()
	elapsed := time.Since(start)

	if elapsed < 15*time.Millisecond {
		t.Fatalf("Wait returned too early: %v", elapsed)
	}
}
