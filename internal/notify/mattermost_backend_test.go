package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestMattermostBackendName(t *testing.T) {
	b := NewMattermostBackend("http://example.com/hook")
	if b.Name() != "mattermost" {
		t.Errorf("expected mattermost, got %s", b.Name())
	}
}

func TestMattermostBackendSendsJSON(t *testing.T) {
	var received map[string]string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected application/json content type")
		}
		json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	b := NewMattermostBackend(ts.URL)
	err := b.Send(Event{Type: "new", Port: 8080, Time: time.Now()})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received["text"] == "" {
		t.Error("expected non-empty text field")
	}
}

func TestMattermostBackendNon2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	b := NewMattermostBackend(ts.URL)
	err := b.Send(Event{Type: "closed", Port: 443, Time: time.Now()})
	if err == nil {
		t.Error("expected error for non-2xx response")
	}
}

func TestMattermostBackendBadURL(t *testing.T) {
	b := NewMattermostBackend("http://127.0.0.1:0/invalid")
	err := b.Send(Event{Type: "new", Port: 22, Time: time.Now()})
	if err == nil {
		t.Error("expected error for bad URL")
	}
}
