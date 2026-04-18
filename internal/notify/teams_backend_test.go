package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTeamsBackendName(t *testing.T) {
	b := NewTeamsBackend("http://example.com")
	if b.Name() != "teams" {
		t.Errorf("expected 'teams', got %q", b.Name())
	}
}

func TestTeamsBackendSendsJSON(t *testing.T) {
	var received map[string]string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("decode error: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	b := NewTeamsBackend(ts.URL)
	err := b.Send(Event{Kind: "new_port", Message: "Port 8080 opened"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received["text"] != "Port 8080 opened" {
		t.Errorf("unexpected text: %q", received["text"])
	}
	if received["title"] != "Portwatch: new_port" {
		t.Errorf("unexpected title: %q", received["title"])
	}
}

func TestTeamsBackendNon2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer ts.Close()

	b := NewTeamsBackend(ts.URL)
	err := b.Send(Event{Kind: "new_port", Message: "test"})
	if err == nil {
		t.Fatal("expected error for non-2xx response")
	}
}

func TestTeamsBackendBadURL(t *testing.T) {
	b := NewTeamsBackend("http://127.0.0.1:0")
	err := b.Send(Event{Kind: "new_port", Message: "test"})
	if err == nil {
		t.Fatal("expected error for bad URL")
	}
}
