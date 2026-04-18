package topports

import (
	"sort"

	"portwatch/internal/scanner"
)

// Entry holds a port and the number of times it has been observed open.
type Entry struct {
	Port  scanner.Port
	Count int
}

// Counter tracks how many times each port has been seen open.
type Counter struct {
	counts map[string]*Entry
}

// New returns an initialised Counter.
func New() *Counter {
	return &Counter{counts: make(map[string]*Entry)}
}

// Record increments the observation count for each port in the slice.
func (c *Counter) Record(ports []scanner.Port) {
	for _, p := range ports {
		key := portKey(p)
		if e, ok := c.counts[key]; ok {
			e.Count++
		} else {
			copy := p
			c.counts[key] = &Entry{Port: copy, Count: 1}
		}
	}
}

// Top returns the n most frequently observed ports, highest count first.
// If n <= 0 all entries are returned.
func (c *Counter) Top(n int) []Entry {
	entries := make([]Entry, 0, len(c.counts))
	for _, e := range c.counts {
		entries = append(entries, *e)
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Count != entries[j].Count {
			return entries[i].Count > entries[j].Count
		}
		return entries[i].Port.Number < entries[j].Port.Number
	})
	if n > 0 && n < len(entries) {
		return entries[:n]
	}
	return entries
}

// Reset clears all counts.
func (c *Counter) Reset() {
	c.counts = make(map[string]*Entry)
}

func portKey(p scanner.Port) string {
	return p.Protocol + ":" + p.String()
}
