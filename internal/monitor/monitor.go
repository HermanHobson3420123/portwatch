package monitor

import (
	"log"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Change represents a detected port state change.
type Change struct {
	Port   scanner.Port
	Opened bool // true if newly opened, false if closed
}

// Monitor watches for port changes at a given interval.
type Monitor struct {
	scanner  *scanner.Scanner
	interval time.Duration
	previous map[string]scanner.Port
	Changes  chan Change
	stop     chan struct{}
}

// New creates a Monitor with the given scanner and poll interval.
func New(s *scanner.Scanner, interval time.Duration) *Monitor {
	return &Monitor{
		scanner:  s,
		interval: interval,
		previous: make(map[string]scanner.Port),
		Changes:  make(chan Change, 32),
		stop:     make(chan struct{}),
	}
}

// Start begins polling in the background. Call Stop to halt.
func (m *Monitor) Start() {
	go m.run()
}

// Stop signals the monitor to cease polling.
func (m *Monitor) Stop() {
	close(m.stop)
}

func (m *Monitor) run() {
	ticker := time.NewTicker(m.interval)
	defer ticker.Stop()

	// Populate baseline on first tick
	if err := m.poll(true); err != nil {
		log.Printf("portwatch monitor: initial scan error: %v", err)
	}

	for {
		select {
		case <-ticker.C:
			if err := m.poll(false); err != nil {
				log.Printf("portwatch monitor: scan error: %v", err)
			}
		case <-m.stop:
			return
		}
	}
}

func (m *Monitor) poll(baseline bool) error {
	ports, err := m.scanner.Scan()
	if err != nil {
		return err
	}

	current := make(map[string]scanner.Port, len(ports))
	for _, p := range ports {
		key := p.String()
		current[key] = p
		if !baseline {
			if _, seen := m.previous[key]; !seen {
				m.Changes <- Change{Port: p, Opened: true}
			}
		}
	}

	if !baseline {
		for key, p := range m.previous {
			if _, stillOpen := current[key]; !stillOpen {
				m.Changes <- Change{Port: p, Opened: false}
			}
		}
	}

	m.previous = current
	return nil
}
