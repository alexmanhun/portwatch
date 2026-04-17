package history

import (
	"fmt"
	"io"
	"sort"
	"time"
)

// Summary holds aggregated statistics for a port history report.
type Summary struct {
	TotalEvents  int
	NewPorts     int
	ClosedPorts  int
	UniquePorts  []int
	FirstSeen    time.Time
	LastSeen     time.Time
}

// Summarize computes a Summary from a slice of Records.
func Summarize(records []Record) Summary {
	if len(records) == 0 {
		return Summary{}
	}

	portSet := map[int]struct{}{}
	s := Summary{
		FirstSeen: records[0].Timestamp,
		LastSeen:  records[0].Timestamp,
	}

	for _, r := range records {
		s.TotalEvents++
		portSet[r.Port] = struct{}{}
		if r.Event == "new" {
			s.NewPorts++
		} else if r.Event == "closed" {
			s.ClosedPorts++
		}
		if r.Timestamp.Before(s.FirstSeen) {
			s.FirstSeen = r.Timestamp
		}
		if r.Timestamp.After(s.LastSeen) {
			s.LastSeen = r.Timestamp
		}
	}

	for p := range portSet {
		s.UniquePorts = append(s.UniquePorts, p)
	}
	sort.Ints(s.UniquePorts)
	return s
}

// PrintReport writes a human-readable summary report to w.
func PrintReport(w io.Writer, s Summary) {
	fmt.Fprintln(w, "=== Port Watch Report ===")
	fmt.Fprintf(w, "Total Events : %d\n", s.TotalEvents)
	fmt.Fprintf(w, "New Ports    : %d\n", s.NewPorts)
	fmt.Fprintf(w, "Closed Ports : %d\n", s.ClosedPorts)
	fmt.Fprintf(w, "Unique Ports : %v\n", s.UniquePorts)
	if !s.FirstSeen.IsZero() {
		fmt.Fprintf(w, "First Event  : %s\n", s.FirstSeen.Format(time.RFC3339))
		fmt.Fprintf(w, "Last Event   : %s\n", s.LastSeen.Format(time.RFC3339))
	}
}
