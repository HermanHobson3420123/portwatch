package filter_test

import (
	"testing"

	"github.com/user/portwatch/internal/filter"
	"github.com/user/portwatch/internal/scanner"
)

func port(number uint16, proto string) scanner.Port {
	return scanner.Port{Number: number, Protocol: proto}
}

func TestIgnored_MatchesExactRule(t *testing.T) {
	f := filter.New([]filter.Rule{{Port: 22, Protocol: "tcp"}})
	if !f.Ignored(port(22, "tcp")) {
		t.Fatal("expected port 22/tcp to be ignored")
	}
}

func TestIgnored_ProtocolMismatch(t *testing.T) {
	f := filter.New([]filter.Rule{{Port: 22, Protocol: "tcp"}})
	if f.Ignored(port(22, "udp")) {
		t.Fatal("expected port 22/udp NOT to be ignored")
	}
}

func TestIgnored_EmptyProtocolMatchesBoth(t *testing.T) {
	f := filter.New([]filter.Rule{{Port: 53, Protocol: ""}})
	if !f.Ignored(port(53, "tcp")) {
		t.Fatal("expected 53/tcp to be ignored")
	}
	if !f.Ignored(port(53, "udp")) {
		t.Fatal("expected 53/udp to be ignored")
	}
}

func TestApply_FiltersCorrectly(t *testing.T) {
	f := filter.New([]filter.Rule{
		{Port: 22, Protocol: "tcp"},
		{Port: 53, Protocol: ""},
	})
	input := []scanner.Port{
		port(22, "tcp"),
		port(53, "tcp"),
		port(53, "udp"),
		port(80, "tcp"),
		port(443, "tcp"),
	}
	result := f.Apply(input)
	if len(result) != 2 {
		t.Fatalf("expected 2 ports after filter, got %d", len(result))
	}
	if result[0].Number != 80 || result[1].Number != 443 {
		t.Fatalf("unexpected ports: %v", result)
	}
}

func TestApply_NoRulesReturnsAll(t *testing.T) {
	f := filter.New(nil)
	input := []scanner.Port{port(80, "tcp"), port(443, "tcp")}
	result := f.Apply(input)
	if len(result) != len(input) {
		t.Fatalf("expected all ports returned, got %d", len(result))
	}
}

func TestFromConfig_ValidEntries(t *testing.T) {
	f, err := filter.FromConfig([]string{"22/tcp", "53", "8080/udp"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !f.Ignored(port(22, "tcp")) {
		t.Error("expected 22/tcp ignored")
	}
	if !f.Ignored(port(53, "tcp")) {
		t.Error("expected 53/tcp ignored")
	}
	if !f.Ignored(port(8080, "udp")) {
		t.Error("expected 8080/udp ignored")
	}
}

func TestFromConfig_InvalidPort(t *testing.T) {
	_, err := filter.FromConfig([]string{"notaport/tcp"})
	if err == nil {
		t.Fatal("expected error for invalid port")
	}
}

func TestFromConfig_InvalidProtocol(t *testing.T) {
	_, err := filter.FromConfig([]string{"80/icmp"})
	if err == nil {
		t.Fatal("expected error for unknown protocol")
	}
}
