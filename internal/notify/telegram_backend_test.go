package notify

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTelegramBackendName(t *testing.T) {
	b := NewTelegramBackend("token", "123")
	if b.Name() != "telegram" {
		t.Fatalf("expected telegram, got %s", b.Name())
	}
}

func TestTelegramBackendSendsJSON(t *testing.T) {
	var got map[string]string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &got)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	b := NewTelegramBackend("mytoken", "456")
	b.apiBase = ts.URL

	err := b.Send(Event{Type: "new_port", Port: 8080})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["chat_id"] != "456" {
		t.Errorf("expected chat_id 456, got %s", got["chat_id"])
	}
	if got["text"] == "" {
		t.Error("expected non-empty text")
	}
}

func TestTelegramBackendNon2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	b := NewTelegramBackend("tok", "789")
	b.apiBase = ts.URL

	if err := b.Send(Event{Type: "closed_port", Port: 22}); err == nil {
		t.Fatal("expected error for non-2xx response")
	}
}

func TestTelegramBackendBadURL(t *testing.T) {
	b := NewTelegramBackend("tok", "1")
	b.apiBase = "http://127.0.0.1:0"
	if err := b.Send(Event{Type: "new_port", Port: 9090}); err == nil {
		t.Fatal("expected error for bad URL")
	}
}
