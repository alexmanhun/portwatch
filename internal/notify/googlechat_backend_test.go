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
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	b := NewGoogleChatBackend(ts.URL)
	ev := alert.Event{Type: alert.NewPort, Port: 8080}
	if err := b.Send(ev); err != nil {
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
	ev := alert.Event{Type: alert.NewPort, Port: 9090}
	if err := b.Send(ev); err == nil {
		t.Fatal("expected error for non-2xx status")
	}
}

func TestGoogleChatBackendBadURL(t *testing.T) {
	b := NewGoogleChatBackend("://bad-url")
	ev := alert.Event{Type: alert.NewPort, Port: 1234}
	if err := b.Send(ev); err == nil {
		t.Fatal("expected error for bad URL")
	}
}
