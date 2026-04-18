package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDiscordBackendName(t *testing.T) {
	d := NewDiscordBackend("https://example.com/hook")
	if d.Name() != "discord" {
		t.Fatalf("expected discord, got %s", d.Name())
	}
}

func TestDiscordBackendSendsJSON(t *testing.T) {
	var received map[string]string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected application/json content-type")
		}
		json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	d := NewDiscordBackend(ts.URL)
	err := d.Send(Event{Type: "new", Port: 8080, Message: "port opened"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received["content"] == "" {
		t.Error("expected non-empty content field")
	}
}

func TestDiscordBackendNon2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	d := NewDiscordBackend(ts.URL)
	err := d.Send(Event{Type: "new", Port: 9090, Message: "test"})
	if err == nil {
		t.Fatal("expected error for non-2xx response")
	}
}

func TestDiscordBackendBadURL(t *testing.T) {
	d := NewDiscordBackend("http://127.0.0.1:0/nope")
	err := d.Send(Event{Type: "closed", Port: 22, Message: "gone"})
	if err == nil {
		t.Fatal("expected error for bad URL")
	}
}
