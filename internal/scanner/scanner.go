package scanner

import (
	"fmt"
	"net"
	"time"
)

// Port represents an open port with its metadata.
type Port struct {
	Number   int
	Protocol string
	Address  string
}

// String returns a human-readable representation of a Port.
func (p Port) String() string {
	return fmt.Sprintf("%s:%d/%s", p.Address, p.Number, p.Protocol)
}

// Scanner scans for open ports on a given host.
type Scanner struct {
	Host    string
	MinPort int
	MaxPort int
	Timeout time.Duration
}

// New creates a new Scanner with sensible defaults.
func New(host string) *Scanner {
	return &Scanner{
		Host:    host,
		MinPort: 1,
		MaxPort: 65535,
		Timeout: 500 * time.Millisecond,
	}
}

// Scan performs a TCP scan over the configured port range and returns open ports.
func (s *Scanner) Scan() ([]Port, error) {
	if err := s.validate(); err != nil {
		return nil, err
	}

	var open []Port

	for port := s.MinPort; port <= s.MaxPort; port++ {
		address := fmt.Sprintf("%s:%d", s.Host, port)
		conn, err := net.DialTimeout("tcp", address, s.Timeout)
		if err != nil {
			continue
		}
		conn.Close()
		open = append(open, Port{
			Number:   port,
			Protocol: "tcp",
			Address:  s.Host,
		})
	}

	return open, nil
}

// validate checks that the Scanner's configuration is valid before scanning.
func (s *Scanner) validate() error {
	if s.MinPort < 1 || s.MaxPort > 65535 {
		return fmt.Errorf("port range must be between 1 and 65535, got %d-%d", s.MinPort, s.MaxPort)
	}
	if s.MinPort > s.MaxPort {
		return fmt.Errorf("MinPort (%d) must not be greater than MaxPort (%d)", s.MinPort, s.MaxPort)
	}
	if s.Host == "" {
		return fmt.Errorf("host must not be empty")
	}
	return nil
}
