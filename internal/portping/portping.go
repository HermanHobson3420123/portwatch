package portping

import (
	"context"
	"fmt"
	"net"
	"time"

	"portwatch/internal/scanner"
)

// Result holds the outcome of a single ping attempt.
type Result struct {
	Port    scanner.Port
	Latency time.Duration
	Alive   bool
	Err     error
}

// Pinger checks whether a port is reachable via TCP/UDP dial.
type Pinger struct {
	timeout time.Duration
}

// New returns a Pinger with the given dial timeout.
func New(timeout time.Duration) *Pinger {
	if timeout <= 0 {
		timeout = 2 * time.Second
	}
	return &Pinger{timeout: timeout}
}

// Ping attempts to connect to the given port and returns a Result.
func (p *Pinger) Ping(ctx context.Context, port scanner.Port) Result {
	addr := fmt.Sprintf("%s:%d", "127.0.0.1", port.Number)
	start := time.Now()

	var d net.Dialer
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	conn, err := d.DialContext(ctx, port.Protocol, addr)
	latency := time.Since(start)
	if err != nil {
		return Result{Port: port, Latency: latency, Alive: false, Err: err}
	}
	_ = conn.Close()
	return Result{Port: port, Latency: latency, Alive: true}
}

// PingAll pings each port and returns all results.
func (p *Pinger) PingAll(ctx context.Context, ports []scanner.Port) []Result {
	results := make([]Result, 0, len(ports))
	for _, port := range ports {
		results = append(results, p.Ping(ctx, port))
	}
	return results
}
