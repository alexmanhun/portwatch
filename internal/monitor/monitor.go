package monitor

import (
	"fmt"
	"time"

	"portwatch/internal/scanner"
)

// PortState holds the last known state of scanned ports.
type PortState struct {
	OpenPorts map[int]bool
	LastScan  time.Time
}

// Monitor watches ports and detects changes over time.
type Monitor struct {
	scanner  *scanner.Scanner
	previous *PortState
	Interval time.Duration
	OnChange func(added, removed []int)
}

// New creates a Monitor with the given scanner and poll interval.
func New(s *scanner.Scanner, interval time.Duration, onChange func(added, removed []int)) *Monitor {
	return &Monitor{
		scanner:  s,
		Interval: interval,
		OnChange: onChange,
	}
}

// Scan performs a single scan and compares with previous state.
func (m *Monitor) Scan(host string, startPort, endPort int) error {
	ports, err := m.scanner.Scan(host, startPort, endPort)
	if err != nil {
		return fmt.Errorf("scan failed: %w", err)
	}

	current := &PortState{
		OpenPorts: make(map[int]bool, len(ports)),
		LastScan:  time.Now(),
	}
	for _, p := range ports {
		current.OpenPorts[p] = true
	}

	if m.previous != nil && m.OnChange != nil {
		added := diff(current.OpenPorts, m.previous.OpenPorts)
		removed := diff(m.previous.OpenPorts, current.OpenPorts)
		if len(added) > 0 || len(removed) > 0 {
			m.OnChange(added, removed)
		}
	}

	m.previous = current
	return nil
}

// diff returns keys present in a but not in b.
func diff(a, b map[int]bool) []int {
	var result []int
	for k := range a {
		if !b[k] {
			result = append(result, k)
		}
	}
	return result
}
