package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
)

func TestCustomEventBackendName(t *testing.T) {
	b := NewCustomEventBackend("http://example.com", nil)
	if b.Name() != "customevent" {
		t.Fatalf("expected customevent, got %s", b.Name())
	}
}

func TestCustomEventBackendSendsJSON(t *testing.T) {
	var received map[string]interface{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	b := NewCustomEventBackend(ts.URL, map[string]string{"X-Token": "abc123"})
	ev := alert.Event{Type: "new_port", Port: 8080, Proto: "tcp", Time: time.Now()}
	if err := b.Send(ev); err != nil {
		t.Fatalf("Send: %v", err)
	}
	if received["event"] != "new_port" {
		t.Errorf("expected event=new_port, got %v", received["event"])
	}
	if int(received["port"].(float64)) != 8080 {
		t.Errorf("expected port=8080, got %v", received["port"])
	}
	if received["proto"] != "tcp" {
		t.Errorf("expected proto=tcp, got %v", received["proto"])
	}
	if _, ok := received["time"]; !ok {
		t.Error("expected time field in payload")
	}
}

func TestCustomEventBackendSendsCustomHeader(t *testing.T) {
	var gotHeader string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotHeader = r.Header.Get("X-Token")
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	b := NewCustomEventBackend(ts.URL, map[string]string{"X-Token": "secret"})
	_ = b.Send(alert.Event{Type: "closed_port", Port: 22, Proto: "tcp", Time: time.Now()})
	if gotHeader != "secret" {
		t.Errorf("expected header secret, got %q", gotHeader)
	}
}

func TestCustomEventBackendNon2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	b := NewCustomEventBackend(ts.URL, nil)
	err := b.Send(alert.Event{Type: "new_port", Port: 9090, Proto: "tcp", Time: time.Now()})
	if err == nil {
		t.Fatal("expected error for non-2xx response")
	}
}

func TestCustomEventBackendBadURL(t *testing.T) {
	b := NewCustomEventBackend("://bad-url", nil)
	err := b.Send(alert.Event{Type: "new_port", Port: 1234, Proto: "tcp", Time: time.Now()})
	if err == nil {
		t.Fatal("expected error for bad URL")
	}
}
