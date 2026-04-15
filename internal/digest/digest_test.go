package digest_test

import (
	"fmt"
	"testing"

	"portwatch/internal/digest"
	"portwatch/internal/scanner"
)

func p(proto string, num uint16) scanner.Port {
	return scanner.Port{Protocol: proto, Number: num}
}

func TestSumIsStable(t *testing.T) {
	list := []scanner.Port{p("tcp", 80), p("tcp", 443)}
	a := digest.Sum(list)
	b := digest.Sum(list)
	if a != b {
		t.Fatalf("expected stable digest, got %q vs %q", a, b)
	}
}

func TestSumIsOrderIndependent(t *testing.T) {
	a := digest.Sum([]scanner.Port{p("tcp", 80), p("tcp", 443)})
	b := digest.Sum([]scanner.Port{p("tcp", 443), p("tcp", 80)})
	if a != b {
		t.Fatalf("order should not affect digest: %q vs %q", a, b)
	}
}

func TestSumDiffersOnChange(t *testing.T) {
	a := digest.Sum([]scanner.Port{p("tcp", 80)})
	b := digest.Sum([]scanner.Port{p("tcp", 8080)})
	if a == b {
		t.Fatal("different ports must produce different digests")
	}
}

func TestSumEmptyList(t *testing.T) {
	s := digest.Sum(nil)
	if s == "" {
		t.Fatal("digest of empty list must not be empty string")
	}
	_ = fmt.Sprintf("digest: %s", s) // ensure fmt import used
}

func TestEqualSamePorts(t *testing.T) {
	list := []scanner.Port{p("udp", 53), p("tcp", 22)}
	if !digest.Equal(list, list) {
		t.Fatal("Equal should return true for identical lists")
	}
}

func TestEqualDifferentPorts(t *testing.T) {
	a := []scanner.Port{p("tcp", 80)}
	b := []scanner.Port{p("tcp", 443)}
	if digest.Equal(a, b) {
		t.Fatal("Equal should return false for different lists")
	}
}

func TestProtocolDistinguished(t *testing.T) {
	a := digest.Sum([]scanner.Port{p("tcp", 53)})
	b := digest.Sum([]scanner.Port{p("udp", 53)})
	if a == b {
		t.Fatal("tcp:53 and udp:53 must produce different digests")
	}
}
