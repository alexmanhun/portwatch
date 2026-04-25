package notify

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/alert"
)

func TestAppriseBackendName(t *testing.T) {
	b := NewAppriseBackend("http://localhost:8000", nil)
	if b.Name() != "apprise" {
		t.Fatalf("expected name 'apprise', got %q", b.Name())
	}
}

func TestAppriseBackendSendsJSON(t *testing.T) {
	var received map[string]interface{}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/notify" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &received)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	b := NewAppriseBackend(ts.URL, []string{"slack://token/channel"})
	event := alert.Event{Port: 8080, Type: alert.NewPort}

	if err := b.Send(event); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if received["title"] == nil {
		t.Fatal("expected title field in payload")
	}
	urls, ok := received["urls"].([]interface{})
	if !ok || len(urls) == 0 {
		t.Fatal("expected non-empty urls field")
	}
	if urls[0] != "slack://token/channel" {
		t.Errorf("unexpected url: %v", urls[0])
	}
}

func TestAppriseBackendNon2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer ts.Close()

	b := NewAppriseBackend(ts.URL, nil)
	err := b.Send(alert.Event{Port: 443, Type: alert.ClosedPort})
	if err == nil {
		t.Fatal("expected error for non-2xx response")
	}
}

func TestAppriseBackendBadURL(t *testing.T) {
	b := NewAppriseBackend("http://127.0.0.1:0", nil)
	err := b.Send(alert.Event{Port: 22, Type: alert.NewPort})
	if err == nil {
		t.Fatal("expected error for unreachable server")
	}
}
