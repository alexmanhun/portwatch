package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGotifyBackendName(t *testing.T) {
	g := NewGotifyBackend("http://localhost", "token")
	if g.Name() != "gotify" {
		t.Errorf("expected gotify, got %s", g.Name())
	}
}

func TestGotifyBackendSendsJSON(t *testing.T) {
	var got map[string]interface{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Gotify-Key") != "mytoken" {
			t.Error("missing token header")
		}
		json.NewDecoder(r.Body).Decode(&got)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	g := NewGotifyBackend(ts.URL, "mytoken")
	err := g.Send(Event{Type: "new_port", Message: "port 8080 opened"})
	if err != nil {
		t.Fatal(err)
	}
	if got["message"] != "port 8080 opened" {
		t.Errorf("unexpected message: %v", got["message"])
	}
}

func TestGotifyBackendNon2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer ts.Close()

	g := NewGotifyBackend(ts.URL, "bad")
	if err := g.Send(Event{Type: "new_port", Message: "x"}); err == nil {
		t.Error("expected error for non-2xx")
	}
}

func TestGotifyBackendBadURL(t *testing.T) {
	g := NewGotifyBackend("http://127.0.0.1:0", "tok")
	if err := g.Send(Event{Type: "new_port", Message: "x"}); err == nil {
		t.Error("expected connection error")
	}
}
