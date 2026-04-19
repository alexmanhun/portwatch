package notify

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRocketChatBackendName(t *testing.T) {
	b := NewRocketChatBackend("http://example.com/hook")
	if b.Name() != "rocketchat" {
		t.Fatalf("expected rocketchat, got %s", b.Name())
	}
}

func TestRocketChatBackendSendsJSON(t *testing.T) {
	var received map[string]string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &received)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	b := NewRocketChatBackend(ts.URL)
	err := b.Send(Event{Type: "new", Port: 8080})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received["text"] == "" {
		t.Fatal("expected non-empty text field")
	}
}

func TestRocketChatBackendNon2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	b := NewRocketChatBackend(ts.URL)
	err := b.Send(Event{Type: "new", Port: 9090})
	if err == nil {
		t.Fatal("expected error on non-2xx response")
	}
}

func TestRocketChatBackendBadURL(t *testing.T) {
	b := NewRocketChatBackend("http://127.0.0.1:0/nope")
	err := b.Send(Event{Type: "closed", Port: 443})
	if err == nil {
		t.Fatal("expected error for bad URL")
	}
}
