package notify

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/alert"
)

func TestBearyChatBackendName(t *testing.T) {
	b := NewBearyChatBackend("http://example.com/hook")
	if b.Name() != "bearychat" {
		t.Fatalf("expected name 'bearychat', got %q", b.Name())
	}
}

func TestBearyChatBackendSendsJSON(t *testing.T) {
	var received bearyChatPayload

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("expected application/json, got %s", ct)
		}
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &received)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	b := NewBearyChatBackend(ts.URL)
	event := alert.Event{Port: 8080, Type: alert.NewPort}

	if err := b.Send(event); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if received.Text == "" {
		t.Error("expected non-empty text in payload")
	}
	if received.Notification == "" {
		t.Error("expected non-empty notification in payload")
	}
}

func TestBearyChatBackendNon2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	b := NewBearyChatBackend(ts.URL)
	event := alert.Event{Port: 443, Type: alert.ClosedPort}

	if err := b.Send(event); err == nil {
		t.Fatal("expected error for non-2xx status, got nil")
	}
}

func TestBearyChatBackendBadURL(t *testing.T) {
	b := NewBearyChatBackend("http://127.0.0.1:0/no-such-endpoint")
	event := alert.Event{Port: 22, Type: alert.NewPort}

	if err := b.Send(event); err == nil {
		t.Fatal("expected error for unreachable URL, got nil")
	}
}
