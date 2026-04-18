package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSMSBackendName(t *testing.T) {
	b := NewSMSBackend("http://example.com", "key", "+1000", "+2000")
	if b.Name() != "sms" {
		t.Fatalf("expected sms, got %s", b.Name())
	}
}

func TestSMSBackendSendsJSON(t *testing.T) {
	var received map[string]string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer testkey" {
			t.Error("missing or wrong Authorization header")
		}
		json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	b := NewSMSBackend(ts.URL, "testkey", "+1000", "+2000")
	err := b.Send(Event{Type: "new_port", Port: 8080})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received["to"] != "+2000" {
		t.Errorf("expected to=+2000, got %s", received["to"])
	}
	if received["from"] != "+1000" {
		t.Errorf("expected from=+1000, got %s", received["from"])
	}
}

func TestSMSBackendNon2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	b := NewSMSBackend(ts.URL, "key", "+1000", "+2000")
	err := b.Send(Event{Type: "new_port", Port: 9090})
	if err == nil {
		t.Fatal("expected error for non-2xx response")
	}
}

func TestSMSBackendBadURL(t *testing.T) {
	b := NewSMSBackend("http://127.0.0.1:0", "key", "+1", "+2")
	err := b.Send(Event{Type: "closed_port", Port: 22})
	if err == nil {
		t.Fatal("expected error for unreachable URL")
	}
}
