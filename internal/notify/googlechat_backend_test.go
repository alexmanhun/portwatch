package notify

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/alert"
)

func TestGoogleChatBackendName(t *testing.T) {
	b := NewGoogleChatBackend("http://example.com/webhook")
	if b.Name() != "googlechat" {
		t.Errorf("expected 'googlechat', got %q", b.Name())
	}
}

func TestGoogleChatBackendSendsJSON(t *testing.T) {
	var received map[string]string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &received)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	b := NewGoogleChatBackend(ts.URL)
	event := alert.Event{Type: alert.NewPortEvent, Port: 8080}
	if err := b.Send(event); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received["text"] == "" {
		t.Error("expected non-empty text field")
	}
}

func TestGoogleChatBackendNon2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	b := NewGoogleChatBackend(ts.URL)
	event := alert.Event{Type: alert.NewPortEvent, Port: 443}
	if err := b.Send(event); err == nil {
		t.Error("expected error on non-2xx response")
	}
}

func TestGoogleChatBackendBadURL(t *testing.T) {
	b := NewGoogleChatBackend("http://127.0.0.1:0/invalid")
	event := alert.Event{Type: alert.ClosedPortEvent, Port: 22}
	if err := b.Send(event); err == nil {
		t.Error("expected error for bad URL")
	}
}
