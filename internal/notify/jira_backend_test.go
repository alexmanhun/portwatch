package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
)

func TestJiraBackendName(t *testing.T) {
	b := NewJiraBackend("https://example.atlassian.net", "OPS", "user@example.com", "token123")
	if b.Name() != "jira" {
		t.Fatalf("expected name 'jira', got %q", b.Name())
	}
}

func TestJiraBackendSendsJSON(t *testing.T) {
	var received map[string]interface{}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/rest/api/2/issue" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("unexpected content-type: %s", ct)
		}
		_ = json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusCreated)
	}))
	defer ts.Close()

	b := NewJiraBackend(ts.URL, "OPS", "user@example.com", "token")
	event := alert.Event{Type: alert.NewPortEvent, Port: 8080, Timestamp: time.Now()}

	if err := b.Send(event); err != nil {
		t.Fatalf("Send returned error: %v", err)
	}

	fields, ok := received["fields"].(map[string]interface{})
	if !ok {
		t.Fatal("missing fields in payload")
	}
	if fields["summary"] == "" {
		t.Error("summary should not be empty")
	}
}

func TestJiraBackendNon2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	b := NewJiraBackend(ts.URL, "OPS", "user", "token")
	event := alert.Event{Type: alert.ClosedPortEvent, Port: 443, Timestamp: time.Now()}

	if err := b.Send(event); err == nil {
		t.Fatal("expected error on non-2xx response")
	}
}

func TestJiraBackendBadURL(t *testing.T) {
	b := NewJiraBackend("://bad-url", "OPS", "user", "token")
	event := alert.Event{Type: alert.NewPortEvent, Port: 22, Timestamp: time.Now()}

	if err := b.Send(event); err == nil {
		t.Fatal("expected error for bad URL")
	}
}
