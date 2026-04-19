package notify

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRedisBackendName(t *testing.T) {
	b := NewRedisBackend("http://localhost", "alerts", "")
	if b.Name() != "redis" {
		t.Fatalf("expected redis, got %s", b.Name())
	}
}

func TestRedisBackendSendsJSON(t *testing.T) {
	var got map[string]interface{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &got)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	b := NewRedisBackend(ts.URL, "portwatch", "secret")
	if err := b.Send("new_port", 8080); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["channel"] != "portwatch" {
		t.Errorf("expected channel portwatch, got %v", got["channel"])
	}
	if got["message"] == "" {
		t.Error("expected non-empty message")
	}
}

func TestRedisBackendNon2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	b := NewRedisBackend(ts.URL, "alerts", "")
	if err := b.Send("closed_port", 9090); err == nil {
		t.Fatal("expected error on non-2xx response")
	}
}

func TestRedisBackendBadURL(t *testing.T) {
	b := NewRedisBackend("://bad-url", "alerts", "")
	if err := b.Send("new_port", 22); err == nil {
		t.Fatal("expected error for bad URL")
	}
}
