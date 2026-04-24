package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestZulipBackendName(t *testing.T) {
	b := NewZulipBackend("https://example.zulipchat.com", "bot@example.com", "key", "alerts", "portwatch")
	if b.Name() != "zulip" {
		t.Fatalf("expected 'zulip', got %q", b.Name())
	}
}

func TestZulipBackendSendsForm(t *testing.T) {
	var captured url.Values
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Fatalf("parse form: %v", err)
		}
		captured = r.Form
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"result": "success", "id": 1})
	}))
	defer ts.Close()

	b := NewZulipBackend(ts.URL, "bot@example.com", "secret", "ops", "alerts")
	event := Event{Type: "new_port", Port: 8080, Detail: "tcp"}

	if err := b.Send(event); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if captured.Get("type") != "stream" {
		t.Errorf("expected type=stream, got %q", captured.Get("type"))
	}
	if captured.Get("to") != "ops" {
		t.Errorf("expected to=ops, got %q", captured.Get("to"))
	}
	if captured.Get("topic") != "alerts" {
		t.Errorf("expected topic=alerts, got %q", captured.Get("topic"))
	}
	msg := captured.Get("content")
	if msg == "" {
		t.Error("expected non-empty content")
	}
}

func TestZulipBackendSendsFormContainsPort(t *testing.T) {
	var captured url.Values
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Fatalf("parse form: %v", err)
		}
		captured = r.Form
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"result": "success", "id": 1})
	}))
	defer ts.Close()

	b := NewZulipBackend(ts.URL, "bot@example.com", "secret", "ops", "alerts")
	event := Event{Type: "new_port", Port: 9090, Detail: "tcp"}

	if err := b.Send(event); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	msg := captured.Get("content")
	if !strings.Contains(msg, "9090") {
		t.Errorf("expected content to contain port 9090, got %q", msg)
	}
}

func TestZulipBackendNon2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer ts.Close()

	b := NewZulipBackend(ts.URL, "bot@example.com", "bad", "ops", "alerts")
	err := b.Send(Event{Type: "new_port", Port: 443})
	if err == nil {
		t.Fatal("expected error for non-2xx response")
	}
}

func TestZulipBackendBadURL(t *testing.T) {
	b := NewZulipBackend("://invalid", "bot@example.com", "key", "ops", "alerts")
	err := b.Send(Event{Type: "new_port", Port: 22})
	if err == nil {
		t.Fatal("expected error for invalid URL")
	}
}
