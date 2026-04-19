package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestZendutyBackendName(t *testing.T) {
	b := NewZendutyBackend("key", "svc", "esc")
	if b.Name() != "zenduty" {
		t.Fatalf("expected zenduty, got %s", b.Name())
	}
}

func TestZendutyBackendSendsJSON(t *testing.T) {
	var got map[string]interface{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Token testkey" {
			t.Errorf("missing or wrong auth header")
		}
		json.NewDecoder(r.Body).Decode(&got)
		w.WriteHeader(http.StatusCreated)
	}))
	defer ts.Close()

	b := NewZendutyBackend("testkey", "svc1", "esc1")
	b.client = ts.Client()
	// Override URL via a round-tripper would be complex; test via direct HTTP.
	// Use a simple approach: swap out the URL by pointing client at test server.
	b.client = &http.Client{
		Transport: rewriteTransport(ts.URL),
	}

	err := b.Send(Event{Type: "new_port", Port: 8080, Message: "port opened"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["service"] != "svc1" {
		t.Errorf("expected service svc1, got %v", got["service"])
	}
}

func TestZendutyBackendNon2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	b := NewZendutyBackend("k", "s", "e")
	b.client = &http.Client{Transport: rewriteTransport(ts.URL)}

	if err := b.Send(Event{Type: "new_port", Port: 22}); err == nil {
		t.Fatal("expected error on 403")
	}
}

func TestZendutyBackendBadURL(t *testing.T) {
	b := NewZendutyBackend("k", "s", "e")
	b.client = &http.Client{Transport: rewriteTransport("http://127.0.0.1:1")}
	if err := b.Send(Event{Type: "closed_port", Port: 443}); err == nil {
		t.Fatal("expected error for unreachable host")
	}
}
