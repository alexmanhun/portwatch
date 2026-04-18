package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestVictorOpsBackendName(t *testing.T) {
	b := NewVictorOpsBackend("http://example.com")
	if b.Name() != "victorops" {
		t.Errorf("expected victorops, got %s", b.Name())
	}
}

func TestVictorOpsBackendSendsJSON(t *testing.T) {
	var received map[string]any
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(200)
	}))
	defer ts.Close()

	b := NewVictorOpsBackend(ts.URL)
	err := b.Send(Event{Port: 8080, Kind: "new_port", Message: "port opened"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received["message_type"] != "WARNING" {
		t.Errorf("expected WARNING, got %v", received["message_type"])
	}
	if received["kind"] != "new_port" {
		t.Errorf("expected new_port, got %v", received["kind"])
	}
}

func TestVictorOpsBackendNon2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(503)
	}))
	defer ts.Close()

	b := NewVictorOpsBackend(ts.URL)
	err := b.Send(Event{Port: 9090, Kind: "closed_port", Message: "port closed"})
	if err == nil {
		t.Fatal("expected error for non-2xx response")
	}
}

func TestVictorOpsBackendBadURL(t *testing.T) {
	b := NewVictorOpsBackend("http://127.0.0.1:0")
	err := b.Send(Event{Port: 22, Kind: "new_port", Message: "ssh opened"})
	if err == nil {
		t.Fatal("expected error for bad URL")
	}
}
