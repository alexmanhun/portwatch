package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOpsGenieBackendName(t *testing.T) {
	b := NewOpsGenieBackend("key123")
	if b.Name() != "opsgenie" {
		t.Fatalf("expected opsgenie, got %s", b.Name())
	}
}

func TestOpsGenieBackendSendsJSON(t *testing.T) {
	var received map[string]interface{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "" {
			t.Error("missing Authorization header")
		}
		json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusAccepted)
	}))
	defer ts.Close()

	b := NewOpsGenieBackend("testkey")
	b.apiURL = ts.URL

	err := b.Send(Event{Type: EventNewPort, Port: 9090, Detail: "new port"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received["message"] == nil {
		t.Error("expected message field in payload")
	}
}

func TestOpsGenieBackendNon2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	b := NewOpsGenieBackend("badkey")
	b.apiURL = ts.URL

	err := b.Send(Event{Type: EventNewPort, Port: 22})
	if err == nil {
		t.Fatal("expected error for non-2xx response")
	}
}

func TestOpsGenieBackendBadURL(t *testing.T) {
	b := NewOpsGenieBackend("key")
	b.apiURL = "http://127.0.0.1:0"
	err := b.Send(Event{Type: EventClosedPort, Port: 80})
	if err == nil {
		t.Fatal("expected error for bad URL")
	}
}
