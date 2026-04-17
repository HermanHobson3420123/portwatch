package trend

import (
	"testing"
)

func TestNewDefaultsWindow(t *testing.T) {
	tr := New(0)
	if tr.window != 2 {
		t.Fatalf("expected window 2, got %d", tr.window)
	}
}

func TestStableWhenFewerThanTwoSamples(t *testing.T) {
	tr := New(5)
	if tr.Direction() != Stable {
		t.Fatal("expected Stable with no samples")
	}
	tr.Record(10)
	if tr.Direction() != Stable {
		t.Fatal("expected Stable with one sample")
	}
}

func TestRisingDirection(t *testing.T) {
	tr := New(5)
	tr.Record(5)
	tr.Record(10)
	if tr.Direction() != Rising {
		t.Fatalf("expected Rising, got %s", tr.Direction())
	}
}

func TestFallingDirection(t *testing.T) {
	tr := New(5)
	tr.Record(10)
	tr.Record(5)
	if tr.Direction() != Falling {
		t.Fatalf("expected Falling, got %s", tr.Direction())
	}
}

func TestStableDirection(t *testing.T) {
	tr := New(5)
	tr.Record(7)
	tr.Record(7)
	if tr.Direction() != Stable {
		t.Fatalf("expected Stable, got %s", tr.Direction())
	}
}

func TestWindowPrunesOldSamples(t *testing.T) {
	tr := New(3)
	for i := 1; i <= 5; i++ {
		tr.Record(i)
	}
	samples := tr.Samples()
	if len(samples) != 3 {
		t.Fatalf("expected 3 samples, got %d", len(samples))
	}
	if samples[0].Count != 3 {
		t.Fatalf("expected oldest kept sample count=3, got %d", samples[0].Count)
	}
}

func TestDirectionString(t *testing.T) {
	if Rising.String() != "rising" {
		t.Fatal("wrong string for Rising")
	}
	if Falling.String() != "falling" {
		t.Fatal("wrong string for Falling")
	}
	if Stable.String() != "stable" {
		t.Fatal("wrong string for Stable")
	}
}
