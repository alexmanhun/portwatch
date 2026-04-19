package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAMQPBackendName(t *testing.T) {
	b := NewAMQPBackend("http://localhost:15672", "%2F", "portwatch", "alerts", "guest", "guest")
	if b.Name() != "amqp" {
		t.Fatalf("expected amqp, got %s", b.Name())
	}
}

func TestAMQPBackendSendsJSON(t *testing.T) {
	var got map[string]interface{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Fatalf("decode: %v", err)
		}
		user, pass, ok := r.BasicAuth()
		if !ok || user != "u" || pass != "p" {
			t.Fatalf("bad auth")
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	b := NewAMQPBackend(ts.URL, "%2F", "ex", "rk", "u", "p")
	if err := b.Send(Event{Type: "new", Port: 8080}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["routing_key"] != "rk" {
		t.Fatalf("expected routing_key rk, got %v", got["routing_key"])
	}
}

func TestAMQPBackendNon2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	b := NewAMQPBackend(ts.URL, "%2F", "ex", "rk", "u", "p")
	if err := b.Send(Event{Type: "closed", Port: 22}); err == nil {
		t.Fatal("expected error for non-2xx status")
	}
}

func TestAMQPBackendBadURL(t *testing.T) {
	b := NewAMQPBackend("http://127.0.0.1:1", "%2F", "ex", "rk", "u", "p")
	if err := b.Send(Event{Type: "new", Port: 9090}); err == nil {
		t.Fatal("expected connection error")
	}
}
