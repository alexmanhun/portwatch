package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMatrixBackendName(t *testing.T) {
	b := NewMatrixBackend("https://matrix.org", "token", "!room:matrix.org")
	if b.Name() != "matrix" {
		t.Fatalf("expected matrix, got %s", b.Name())
	}
}

func TestMatrixBackendSendsJSON(t *testing.T) {
	var got map[string]string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer testtoken" {
			t.Error("missing or wrong Authorization header")
		}
		json.NewDecoder(r.Body).Decode(&got)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"event_id":"$abc"}`))
	}))
	defer ts.Close()

	b := NewMatrixBackend(ts.URL, "testtoken", "!room:example.org")
	err := b.Send(Event{Type: "new_port", Port: 8080})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["msgtype"] != "m.text" {
		t.Errorf("expected m.text, got %s", got["msgtype"])
	}
	if got["body"] == "" {
		t.Error("expected non-empty body")
	}
}

func TestMatrixBackendNon2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	b := NewMatrixBackend(ts.URL, "bad", "!room:example.org")
	err := b.Send(Event{Type: "new_port", Port: 22})
	if err == nil {
		t.Fatal("expected error for non-2xx response")
	}
}

func TestMatrixBackendBadURL(t *testing.T) {
	b := NewMatrixBackend("http://127.0.0.1:0", "tok", "!r:x")
	err := b.Send(Event{Type: "closed_port", Port: 443})
	if err == nil {
		t.Fatal("expected error for unreachable server")
	}
}
