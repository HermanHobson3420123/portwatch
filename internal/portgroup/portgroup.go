package portgroup

import (
	"fmt"
	"sort"

	"portwatch/internal/scanner"
)

// Group represents a named collection of ports sharing a common label.
type Group struct {
	Name  string
	Ports []scanner.Port
}

// Grouper buckets ports by a user-supplied key function.
type Grouper struct {
	keyFn func(scanner.Port) string
}

// New returns a Grouper that partitions ports using keyFn.
func New(keyFn func(scanner.Port) string) *Grouper {
	return &Grouper{keyFn: keyFn}
}

// ByProtocol returns a Grouper that groups ports by protocol.
func ByProtocol() *Grouper {
	return New(func(p scanner.Port) string { return p.Protocol })
}

// ByPortRange returns a Grouper that buckets ports into well-known ranges.
func ByPortRange() *Grouper {
	return New(func(p scanner.Port) string {
		switch {
		case p.Number < 1024:
			return "system (0-1023)"
		case p.Number < 49152:
			return "registered (1024-49151)"
		default:
			return "dynamic (49152+)"
		}
	})
}

// Group partitions ports and returns sorted groups.
func (g *Grouper) Group(ports []scanner.Port) []Group {
	buckets := make(map[string][]scanner.Port)
	for _, p := range ports {
		k := g.keyFn(p)
		buckets[k] = append(buckets[k], p)
	}
	groups := make([]Group, 0, len(buckets))
	for name, ps := range buckets {
		sort.Slice(ps, func(i, j int) bool { return ps[i].Number < ps[j].Number })
		groups = append(groups, Group{Name: name, Ports: ps})
	}
	sort.Slice(groups, func(i, j int) bool { return groups[i].Name < groups[j].Name })
	return groups
}

// Summary returns a human-readable one-liner per group.
func Summary(groups []Group) string {
	out := ""
	for _, g := range groups {
		out += fmt.Sprintf("%s: %d port(s)\n", g.Name, len(g.Ports))
	}
	return out
}
