package history

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func baseRecords() []Record {
	now := time.Now()
	return []Record{
		{Port: 80, Event: "new", Timestamp: now.Add(-2 * time.Hour)},
		{Port: 443, Event: "new", Timestamp: now.Add(-1 * time.Hour)},
		{Port: 80, Event: "closed", Timestamp: now},
	}
}

func TestSummarizeEmpty(t *testing.T) {
	s := Summarize(nil)
	if s.TotalEvents != 0 {
		t.Errorf("expected 0 events, got %d", s.TotalEvents)
	}
}

func TestSummarizeCounts(t *testing.T) {
	s := Summarize(baseRecords())
	if s.TotalEvents != 3 {
		t.Errorf("expected 3 total, got %d", s.TotalEvents)
	}
	if s.NewPorts != 2 {
		t.Errorf("expected 2 new, got %d", s.NewPorts)
	}
	if s.ClosedPorts != 1 {
		t.Errorf("expected 1 closed, got %d", s.ClosedPorts)
	}
}

func TestSummarizeUniquePorts(t *testing.T) {
	s := Summarize(baseRecords())
	if len(s.UniquePorts) != 2 {
		t.Fatalf("expected 2 unique ports, got %d", len(s.UniquePorts))
	}
	if s.UniquePorts[0] != 80 || s.UniquePorts[1] != 443 {
		t.Errorf("unexpected ports: %v", s.UniquePorts)
	}
}

func TestSummarizeTimeRange(t *testing.T) {
	records := baseRecords()
	s := Summarize(records)
	if !s.FirstSeen.Before(s.LastSeen) {
		t.Errorf("expected FirstSeen before LastSeen")
	}
}

func TestPrintReport(t *testing.T) {
	s := Summarize(baseRecords())
	var buf bytes.Buffer
	PrintReport(&buf, s)
	out := buf.String()
	for _, want := range []string{"Report", "Total", "New", "Closed", "Unique"} {
		if !strings.Contains(out, want) {
			t.Errorf("report missing %q", want)
		}
	}
}
