package audit

import "time"

// Filter holds optional constraints for querying audit events.
type Filter struct {
	Kind  EventKind
	Since time.Time
	Proto string
}

// Query returns events matching all non-zero fields of f.
func (l *Log) Query(f Filter) []Event {
	l.mu.Lock()
	defer l.mu.Unlock()
	var out []Event
	for _, ev := range l.events {
		if f.Kind != "" && ev.Kind != f.Kind {
			continue
		}
		if !f.Since.IsZero() && ev.Timestamp.Before(f.Since) {
			continue
		}
		if f.Proto != "" && ev.Port.Proto != f.Proto {
			continue
		}
		out = append(out, ev)
	}
	return out
}
