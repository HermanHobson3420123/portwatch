package portdiff

import (
	"fmt"
	"strings"

	"portwatch/internal/scanner"
)

// Summary holds a human-readable diff between two port snapshots.
type Summary struct {
	Opened []scanner.Port
	Closed []scanner.Port
}

// Empty returns true when there are no changes.
func (s Summary) Empty() bool {
	return len(s.Opened) == 0 && len(s.Closed) == 0
}

// String returns a compact multi-line representation of the diff.
func (s Summary) String() string {
	var b strings.Builder
	for _, p := range s.Opened {
		fmt.Fprintf(&b, "+ %s\n", p)
	}
	for _, p := range s.Closed {
		fmt.Fprintf(&b, "- %s\n", p)
	}
	return strings.TrimRight(b.String(), "\n")
}

// Compare returns a Summary describing ports that moved between prev and next.
func Compare(prev, next []scanner.Port) Summary {
	prevSet := toSet(prev)
	nextSet := toSet(next)

	var s Summary
	for k, p := range nextSet {
		if _, ok := prevSet[k]; !ok {
			s.Opened = append(s.Opened, p)
		}
	}
	for k, p := range prevSet {
		if _, ok := nextSet[k]; !ok {
			s.Closed = append(s.Closed, p)
		}
	}
	return s
}

func toSet(ports []scanner.Port) map[string]scanner.Port {
	m := make(map[string]scanner.Port, len(ports))
	for _, p := range ports {
		m[fmt.Sprintf("%s:%d", p.Protocol, p.Number)] = p
	}
	return m
}
