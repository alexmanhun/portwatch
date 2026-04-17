package history

import (
	"testing"
	"time"
)

var now = time.Now()

func filterRecords() []Record {
	return []Record{
		{Port: 80, Event: "opened", Timestamp: now.Add(-3 * time.Hour)},
		{Port: 443, Event: "opened", Timestamp: now.Add(-2 * time.Hour)},
		{Port: 80, Event: "closed", Timestamp: now.Add(-1 * time.Hour)},
		{Port: 8080, Event: "opened", Timestamp: now},
	}
}

func TestFilterByPort(t *testing.T) {
	results := Filter(filterRecords(), FilterOptions{Port: 80})
	if len(results) != 2 {
		t.Fatalf("expected 2, got %d", len(results))
	}
}

func TestFilterByEvent(t *testing.T) {
	results := Filter(filterRecords(), FilterOptions{Event: "closed"})
	if len(results) != 1 || results[0].Port != 80 {
		t.Fatalf("unexpected results: %+v", results)
	}
}

func TestFilterBySince(t *testing.T) {
	results := Filter(filterRecords(), FilterOptions{Since: now.Add(-90 * time.Minute)})
	if len(results) != 2 {
		t.Fatalf("expected 2, got %d", len(results))
	}
}

func TestFilterByUntil(t *testing.T) {
	results := Filter(filterRecords(), FilterOptions{Until: now.Add(-90 * time.Minute)})
	if len(results) != 2 {
		t.Fatalf("expected 2, got %d", len(results))
	}
}

func TestFilterLimit(t *testing.T) {
	results := Filter(filterRecords(), FilterOptions{Limit: 2})
	if len(results) != 2 {
		t.Fatalf("expected 2, got %d", len(results))
	}
	if results[0].Port != 80 && results[0].Event != "closed" {
		t.Fatalf("expected last 2 records")
	}
}

func TestFilterNoOpts(t *testing.T) {
	results := Filter(filterRecords(), FilterOptions{})
	if len(results) != 4 {
		t.Fatalf("expected all 4, got %d", len(results))
	}
}
