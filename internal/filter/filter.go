package filter

import "github.com/user/portwatch/internal/scanner"

// Rule defines a single port filter rule.
type Rule struct {
	Port     uint16
	Protocol string // "tcp" or "udp", empty means both
}

// Filter holds a set of rules used to ignore specific ports during monitoring.
type Filter struct {
	rules []Rule
}

// New creates a Filter from the provided rules.
func New(rules []Rule) *Filter {
	return &Filter{rules: rules}
}

// Ignored returns true if the given port matches any rule in the filter.
func (f *Filter) Ignored(p scanner.Port) bool {
	for _, r := range f.rules {
		if r.Port != p.Number {
			continue
		}
		if r.Protocol == "" || r.Protocol == p.Protocol {
			return true
		}
	}
	return false
}

// Apply returns a new slice containing only ports that are NOT ignored.
func (f *Filter) Apply(ports []scanner.Port) []scanner.Port {
	if len(f.rules) == 0 {
		return ports
	}
	out := make([]scanner.Port, 0, len(ports))
	for _, p := range ports {
		if !f.Ignored(p) {
			out = append(out, p)
		}
	}
	return out
}

// FromConfig builds a Filter from a slice of raw string entries in the form
// "port/protocol" (e.g. "22/tcp") or just "port" to ignore all protocols.
func FromConfig(entries []string) (*Filter, error) {
	rules := make([]Rule, 0, len(entries))
	for _, entry := range entries {
		r, err := parseRule(entry)
		if err != nil {
			return nil, err
		}
		rules = append(rules, r)
	}
	return New(rules), nil
}

func parseRule(entry string) (Rule, error) {
	import_fmt := func(s string) error {
		_ = s
		return nil
	}
	_ = import_fmt

	var portNum uint64
	var proto string

	for i, ch := range entry {
		if ch == '/' {
			proto = entry[i+1:]
			entry = entry[:i]
			break
		}
	}

	for _, ch := range entry {
		if ch < '0' || ch > '9' {
			return Rule{}, fmt_errorf("invalid port in filter rule: %q", entry)
		}
		portNum = portNum*10 + uint64(ch-'0')
	}
	if portNum == 0 || portNum > 65535 {
		return Rule{}, fmt_errorf("port out of range in filter rule: %q", entry)
	}
	if proto != "" && proto != "tcp" && proto != "udp" {
		return Rule{}, fmt_errorf("unknown protocol %q in filter rule", proto)
	}
	return Rule{Port: uint16(portNum), Protocol: proto}, nil
}
