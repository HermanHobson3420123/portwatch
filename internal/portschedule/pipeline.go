package portschedule

import (
	"context"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Gate wraps a Schedule and filters scanner.Port slices to only those
// matching an active scheduled window.
type Gate struct {
	sched *Schedule
	clock func() time.Time
}

// NewGate returns a Gate backed by the given Schedule.
func NewGate(s *Schedule) *Gate {
	return &Gate{sched: s, clock: time.Now}
}

// Filter returns only ports that fall within an active scheduled window.
// If no windows are active, all ports are returned unchanged.
func (g *Gate) Filter(ports []scanner.Port) []scanner.Port {
	active := g.sched.ActiveNow(g.clock())
	if len(active) == 0 {
		return ports
	}
	var out []scanner.Port
	for _, p := range ports {
		if matchesAny(p, active) {
			out = append(out, p)
		}
	}
	return out
}

func matchesAny(p scanner.Port, entries []Entry) bool {
	for _, e := range entries {
		protoMatch := e.Protocol == "" || e.Protocol == p.Proto
		portMatch := p.Number >= e.PortLow && p.Number <= e.PortHigh
		if protoMatch && portMatch {
			return true
		}
	}
	return false
}

// NewPipeline reads ports from in, applies the gate filter, and forwards to out.
func NewPipeline(ctx context.Context, g *Gate, in <-chan []scanner.Port) <-chan []scanner.Port {
	out := make(chan []scanner.Port)
	go func() {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				return
			case ports, ok := <-in:
				if !ok {
					return
				}
				filtered := g.Filter(ports)
				select {
				case out <- filtered:
				case <-ctx.Done():
					return
				}
			}
		}
	}()
	return out
}
