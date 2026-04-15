package ratelimit_test

import (
	"testing"
	"time"

	"portwatch/internal/ratelimit"
)

func TestAllowFirstEventPasses(t *testing.T) {
	l := ratelimit.New(100 * time.Millisecond)
	if !l.Allow("tcp:8080") {
		t.Fatal("expected first event to be allowed")
	}
}

func TestAllowSecondEventSuppressed(t *testing.T) {
	l := ratelimit.New(100 * time.Millisecond)
	l.Allow("tcp:8080")
	if l.Allow("tcp:8080") {
		t.Fatal("expected second event within cooldown to be suppressed")
	}
}

func TestAllowAfterCooldownPasses(t *testing.T) {
	l := ratelimit.New(20 * time.Millisecond)
	l.Allow("tcp:8080")
	time.Sleep(30 * time.Millisecond)
	if !l.Allow("tcp:8080") {
		t.Fatal("expected event after cooldown to be allowed")
	}
}

func TestDifferentKeysAreIndependent(t *testing.T) {
	l := ratelimit.New(100 * time.Millisecond)
	l.Allow("tcp:8080")
	if !l.Allow("udp:9090") {
		t.Fatal("expected different key to be allowed independently")
	}
}

func TestResetAllowsImmediateRetry(t *testing.T) {
	l := ratelimit.New(100 * time.Millisecond)
	l.Allow("tcp:8080")
	l.Reset("tcp:8080")
	if !l.Allow("tcp:8080") {
		t.Fatal("expected event to pass after Reset")
	}
}

func TestFlushClearsAllKeys(t *testing.T) {
	l := ratelimit.New(100 * time.Millisecond)
	l.Allow("tcp:8080")
	l.Allow("udp:53")
	l.Flush()
	if !l.Allow("tcp:8080") {
		t.Fatal("expected tcp:8080 to pass after Flush")
	}
	if !l.Allow("udp:53") {
		t.Fatal("expected udp:53 to pass after Flush")
	}
}
