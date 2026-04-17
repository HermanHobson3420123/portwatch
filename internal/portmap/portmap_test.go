package portmap_test

import (
	"testing"

	"github.com/user/portwatch/internal/portmap"
	"github.com/user/portwatch/internal/scanner"
)

func p(proto string, number int) scanner.Port {
	return scanner.Port{Proto: proto, Number: number}
}

func TestSetAndHas(t *testing.T) {
	m := portmap.New()
	port := p("tcp", 8080)
	if m.Has(port) {
		t.Fatal("expected port to be absent initially")
	}
	m.Set(port)
	if !m.Has(port) {
		t.Fatal("expected port to be present after Set")
	}
}

func TestDelete(t *testing.T) {
	m := portmap.New()
	port := p("tcp", 443)
	m.Set(port)
	m.Delete(port)
	if m.Has(port) {
		t.Fatal("expected port to be absent after Delete")
	}
}

func TestLen(t *testing.T) {
	m := portmap.New()
	if m.Len() != 0 {
		t.Fatalf("expected 0, got %d", m.Len())
	}
	m.Set(p("tcp", 80))
	m.Set(p("udp", 53))
	if m.Len() != 2 {
		t.Fatalf("expected 2, got %d", m.Len())
	}
}

func TestAll(t *testing.T) {
	m := portmap.New()
	m.Set(p("tcp", 22))
	m.Set(p("tcp", 80))
	all := m.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 ports, got %d", len(all))
	}
}

func TestProtoDistinction(t *testing.T) {
	m := portmap.New()
	m.Set(p("tcp", 53))
	m.Set(p("udp", 53))
	if m.Len() != 2 {
		t.Fatalf("tcp and udp on same port should be distinct, got %d", m.Len())
	}
}

func TestSetOverwrites(t *testing.T) {
	m := portmap.New()
	m.Set(p("tcp", 8080))
	m.Set(p("tcp", 8080))
	if m.Len() != 1 {
		t.Fatalf("duplicate Set should not grow map, got %d", m.Len())
	}
}
