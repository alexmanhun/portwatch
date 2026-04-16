package monitor

import (
	"time"

	"portwatch/internal/alert"
	"portwatch/internal/scanner"
)

// Monitor watches a port range and fires callbacks on changes.
type Monitor struct {
	scanner  *scanner.Scanner
	notifier *alert.Notifier
	baseline map[int]bool
	interval time.Duration
}

// New creates a Monitor for the given port range and alert notifier.
func New(start, end int, interval time.Duration, n *alert.Notifier) (*Monitor, error) {
	s, err := scanner.New(start, end)
	if err != nil {
		return nil, err
	}
	return &Monitor{
		scanner:  s,
		notifier: n,
		baseline: make(map[int]bool),
		interval: interval,
	}, nil
}

// Baseline captures the current open ports as the known-good state.
func (m *Monitor) Baseline() error {
	ports, err := m.scanner.Scan()
	if err != nil {
		return err
	}
	m.baseline = make(map[int]bool, len(ports))
	for _, p := range ports {
		m.baseline[p] = true
	}
	return nil
}

// Check scans once and emits alerts for any changes since the last baseline.
func (m *Monitor) Check() error {
	ports, err := m.scanner.Scan()
	if err != nil {
		return err
	}
	current := make(map[int]bool, len(ports))
	for _, p := range ports {
		current[p] = true
	}
	for p := range current {
		if !m.baseline[p] {
			_ = m.notifier.Notify(alert.NewPortEvent(p))
		}
	}
	for p := range m.baseline {
		if !current[p] {
			_ = m.notifier.Notify(alert.ClosedPortEvent(p))
		}
	}
	m.baseline = current
	return nil
}

// Run starts the monitoring loop and blocks until stop is closed.
func (m *Monitor) Run(stop <-chan struct{}) error {
	if err := m.Baseline(); err != nil {
		return err
	}
	ticker := time.NewTicker(m.interval)
	defer ticker.Stop()
	for {
		select {
		case <-stop:
			return nil
		case <-ticker.C:
			_ = m.Check()
		}
	}
}
