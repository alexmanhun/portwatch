package notify

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSlackBackendName(t *testing.T) {
	s := NewSlackBackend("http://example.com/hook")
	if s.Name() != "slack" {
		t.Errorf("expected name 'slack', got %q", s.Name())
	}
}

func TestSlackBackendSendsJSON(t *testing.T) {
	var received map[string]string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &received)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	s := NewSlackBackend(ts.URL)
	err := s.Send(Event{Type: "new", Port: 8080, Message: "port opened"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(received["text"], "8080") {
		t.Errorf("expected port 8080 in message, got: %q", received["text"])
	}
	if !strings.Contains(received["text"], "portwatch alert") {
		t.Errorf("expected 'portwatch alert' in message, got: %q", received["text"])
	}
}

func TestSlackBackendNon2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	s := NewSlackBackend(ts.URL)
	err := s.Send(Event{Type: "closed", Port: 443, Message: "port closed"})
	if err == nil {
		t.Fatal("expected error for non-2xx response")
	}
}

func TestSlackBackendBadURL(t *testing.T) {
	s := NewSlackBackend("http://127.0.0.1:0/no-such-server")
	err := s.Send(Event{Type: "new", Port: 22, Message: "test"})
	if err == nil {
		t.Fatal("expected error for unreachable server")
	}
}
