package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/alert"
)

func TestLineBackendName(t *testing.T) {
	b := NewLineBackend("tok")
	if b.Name() != "line" {
		t.Fatalf("expected name 'line', got %q", b.Name())
	}
}

func TestLineBackendSendsJSON(t *testing.T) {
	var gotAuth, gotContentType string
	var body map[string]string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		gotContentType = r.Header.Get("Content-Type")
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "bad body", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	b := NewLineBackend("mytoken")
	b.endpoint = ts.URL

	ev := alert.Event{Type: "new", Port: 8080, Message: "opened"}
	if err := b.Send(ev); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if gotAuth != "Bearer mytoken" {
		t.Errorf("expected auth header 'Bearer mytoken', got %q", gotAuth)
	}
	if gotContentType != "application/json" {
		t.Errorf("expected content-type 'application/json', got %q", gotContentType)
	}
	if body["message"] == "" {
		t.Error("expected non-empty message in payload")
	}
}

func TestLineBackendNon2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer ts.Close()

	b := NewLineBackend("bad-token")
	b.endpoint = ts.URL

	err := b.Send(alert.Event{Type: "closed", Port: 22, Message: "gone"})
	if err == nil {
		t.Fatal("expected error for non-2xx response")
	}
}

func TestLineBackendBadURL(t *testing.T) {
	b := NewLineBackend("tok")
	b.endpoint = "://invalid"

	err := b.Send(alert.Event{Type: "new", Port: 9090, Message: "up"})
	if err == nil {
		t.Fatal("expected error for invalid URL")
	}
}
