package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSignalRBackendName(t *testing.T) {
	b := NewSignalRBackend("http://example.com", "testhub", "key")
	if b.Name() != "signalr" {
		t.Errorf("expected signalr, got %s", b.Name())
	}
}

func TestSignalRBackendSendsJSON(t *testing.T) {
	var received map[string]interface{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("decode error: %v", err)
		}
		if r.Header.Get("Authorization") != "Bearer testkey" {
			t.Errorf("missing auth header")
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	b := NewSignalRBackend(ts.URL, "myhub", "testkey")
	err := b.Send(Event{Type: "new_port", Port: 8080, Message: "port opened"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received["target"] != "portwatch" {
		t.Errorf("unexpected target: %v", received["target"])
	}
}

func TestSignalRBackendNon2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	b := NewSignalRBackend(ts.URL, "hub", "")
	err := b.Send(Event{Type: "new_port", Port: 9090, Message: "test"})
	if err == nil {
		t.Fatal("expected error for non-2xx status")
	}
}

func TestSignalRBackendBadURL(t *testing.T) {
	b := NewSignalRBackend("://bad-url", "hub", "")
	err := b.Send(Event{Type: "new_port", Port: 1234, Message: "test"})
	if err == nil {
		t.Fatal("expected error for bad URL")
	}
}
