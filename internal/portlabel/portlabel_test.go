package portlabel

import "testing"

func TestLabelKnownPort(t *testing.T) {
	if got := Label(80, "tcp"); got != "http" {
		t.Fatalf("expected http, got %q", got)
	}
}

func TestLabelFallsBackToGeneric(t *testing.T) {
	// port 53 with empty protocol should still match generic entry
	if got := Label(53, ""); got != "dns" {
		t.Fatalf("expected dns, got %q", got)
	}
}

func TestLabelUnknownPort(t *testing.T) {
	if got := Label(9999, "tcp"); got != "" {
		t.Fatalf("expected empty, got %q", got)
	}
}

func TestAnnotateKnownPort(t *testing.T) {
	if got := Annotate(443, "tcp"); got != "https" {
		t.Fatalf("expected https, got %q", got)
	}
}

func TestAnnotateUnknownReturnsUnknown(t *testing.T) {
	if got := Annotate(19999, "tcp"); got != "unknown" {
		t.Fatalf("expected unknown, got %q", got)
	}
}

func TestItoa(t *testing.T) {
	cases := []struct {
		in  uint16
		out string
	}{
		{0, "0"},
		{1, "1"},
		{443, "443"},
		{65535, "65535"},
	}
	for _, c := range cases {
		if got := itoa(c.in); got != c.out {
			t.Errorf("itoa(%d) = %q, want %q", c.in, got, c.out)
		}
	}
}
