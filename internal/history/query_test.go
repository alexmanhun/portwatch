package history

import (
	"testing"
	"time"
)

func seedHistory(t *testing.T) *History {
	t.Helper()
	h, _ := New(t.TempDir() + "/h.json")
	h.Record(80, "opened")
	h.Record(443, "opened")
	h.Record(80, "closed")
	h.Record(8080, "opened")
	return h
}

func TestQueryByPort(t *testing.T) {
	h := seedHistory(t)
	res := h.Query(Filter{Port: 80})
	if len(res) != 2 {
		t.Fatalf("expected 2, got %d", len(res))
	}
}

func TestQueryByEvent(t *testing.T) {
	h := seedHistory(t)
	res := h.Query(Filter{Event: "closed"})
	if len(res) != 1 {
		t.Fatalf("expected 1 closed event, got %d", len(res))
	}
	if res[0].Port != 80 {
		t.Errorf("expected port 80, got %d", res[0].Port)
	}
}

func TestQueryBySince(t *testing.T) {
	h := seedHistory(t)
	future := time.Now().Add(time.Hour)
	res := h.Query(Filter{Since: future})
	if len(res) != 0 {
		t.Errorf("expected 0 results for future since, got %d", len(res))
	}
}

func TestLast(t *testing.T) {
	h := seedHistory(t)
	res := h.Last(2)
	if len(res) != 2 {
		t.Fatalf("expected 2, got %d", len(res))
	}
	if res[1].Port != 8080 {
		t.Errorf("expected last entry port 8080, got %d", res[1].Port)
	}
}

func TestLastAll(t *testing.T) {
	h := seedHistory(t)
	res := h.Last(0)
	if len(res) != 4 {
		t.Fatalf("expected all 4, got %d", len(res))
	}
}
