package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPagerDutyBackendName(t *testing.T) {
	b := NewPagerDutyBackend("key123")
	if b.Name() != "pagerduty" {
		t.Errorf("expected pagerduty, got %s", b.Name())
	}
}

func TestPagerDutyBackendSendsJSON(t *testing.T) {
	var received map[string]interface{}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("decode error: %v", err)
		}
		w.WriteHeader(http.StatusAccepted)
	}))
	defer ts.Close()

	b := NewPagerDutyBackend("testkey")
	b.url = ts.URL

	err := b.Send(Event{Port: 8080, Type: "new", Message: "port 8080 opened"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if received["routing_key"] != "testkey" {
		t.Errorf("expected routing_key testkey, got %v", received["routing_key"])
	}
	payload, ok := received["payload"].(map[string]interface{})
	if !ok {
		t.Fatal("missing payload")
	}
	if payload["summary"] != "port 8080 opened" {
		t.Errorf("unexpected summary: %v", payload["summary"])
	}
}

func TestPagerDutyBackendNon2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer ts.Close()

	b := NewPagerDutyBackend("key")
	b.url = ts.URL

	err := b.Send(Event{Port: 443, Type: "closed", Message: "port 443 closed"})
	if err == nil {
		t.Error("expected error for non-2xx response")
	}
}
