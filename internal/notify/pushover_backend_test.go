package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPushoverBackendName(t *testing.T) {
	b := NewPushoverBackend("tok", "usr")
	if b.Name() != "pushover" {
		t.Fatalf("expected pushover, got %s", b.Name())
	}
}

func TestPushoverBackendSendsJSON(t *testing.T) {
	var got map[string]string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&got)
		w.WriteHeader(200)
	}))
	defer ts.Close()

	b := NewPushoverBackend("mytoken", "myuser")
	b.apiURL = ts.URL

	err := b.Send(Event{Type: "new_port", Message: "port 8080 opened"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["token"] != "mytoken" {
		t.Errorf("expected token mytoken, got %s", got["token"])
	}
	if got["user"] != "myuser" {
		t.Errorf("expected user myuser, got %s", got["user"])
	}
	if got["message"] != "port 8080 opened" {
		t.Errorf("unexpected message: %s", got["message"])
	}
}

func TestPushoverBackendNon2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(429)
	}))
	defer ts.Close()

	b := NewPushoverBackend("tok", "usr")
	b.apiURL = ts.URL

	if err := b.Send(Event{Type: "new_port", Message: "x"}); err == nil {
		t.Fatal("expected error for non-2xx")
	}
}

func TestPushoverBackendBadURL(t *testing.T) {
	b := NewPushoverBackend("tok", "usr")
	b.apiURL = "http://127.0.0.1:0"
	if err := b.Send(Event{Type: "new_port", Message: "x"}); err == nil {
		t.Fatal("expected error for bad URL")
	}
}
