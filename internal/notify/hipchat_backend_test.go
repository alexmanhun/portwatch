package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHipChatBackendName(t *testing.T) {
	b := NewHipChatBackend("tok", "123", "")
	if b.Name() != "hipchat" {
		t.Fatalf("expected hipchat, got %s", b.Name())
	}
}

func TestHipChatBackendSendsJSON(t *testing.T) {
	var got map[string]interface{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Fatal(err)
		}
		w.WriteHeader(204)
	}))
	defer ts.Close()

	b := NewHipChatBackend("mytoken", "42", ts.URL)
	if err := b.Send(Event{Type: "new", Port: 8080}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["color"] != "yellow" {
		t.Errorf("expected yellow color for new event, got %v", got["color"])
	}
}

func TestHipChatBackendNon2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(403)
	}))
	defer ts.Close()

	b := NewHipChatBackend("tok", "1", ts.URL)
	if err := b.Send(Event{Type: "closed", Port: 22}); err == nil {
		t.Fatal("expected error for non-2xx response")
	}
}

func TestHipChatBackendBadURL(t *testing.T) {
	b := NewHipChatBackend("tok", "1", "http://127.0.0.1:0")
	if err := b.Send(Event{Type: "new", Port: 9090}); err == nil {
		t.Fatal("expected error for bad URL")
	}
}

func TestColorForEvent(t *testing.T) {
	if colorForEvent("new") != "yellow" {
		t.Error("new should be yellow")
	}
	if colorForEvent("closed") != "red" {
		t.Error("closed should be red")
	}
	if colorForEvent("unknown") != "gray" {
		t.Error("unknown should be gray")
	}
}
