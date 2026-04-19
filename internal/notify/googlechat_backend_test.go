package notify

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGoogleChatBackendName(t *testing.T) {
	b := NewGoogleChatBackend("http://example.com")
	if b.Name() != "googlechat" {
		t.Fatalf("expected googlechat, got %s", b.Name())
	}
}

func TestGoogleChatBackendSendsJSON(t *testing.T) {
	var received map[string]string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &received)
		w.WriteHeader(200)
	}))
	defer ts.Close()

	b := NewGoogleChatBackend(ts.URL)
	err := b.Send(Event{Type: "new", Port: 8080, Message: "opened"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received["text"] == "" {
		t.Fatal("expected text field in payload")
	}
}

func TestGoogleChatBackendNon2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
	}))
	defer ts.Close()

	b := NewGoogleChatBackend(ts.URL)
	err := b.Send(Event{Type: "new", Port: 9090, Message: "opened"})
	if err == nil {
		t.Fatal("expected error on non-2xx response")
	}
}

func TestGoogleChatBackendBadURL(t *testing.T) {
	b := NewGoogleChatBackend("http://127.0.0.1:0")
	err := b.Send(Event{Type: "closed", Port: 22, Message: "closed"})
	if err == nil {
		t.Fatal("expected error on bad URL")
	}
}
