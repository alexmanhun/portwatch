package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/alert"
)

func TestPushbulletBackendName(t *testing.T) {
	b := NewPushbulletBackend("key")
	if b.Name() != "pushbullet" {
		t.Fatalf("expected pushbullet, got %s", b.Name())
	}
}

func TestPushbulletBackendSendsJSON(t *testing.T) {
	var got map[string]interface{}
	var gotToken string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotToken = r.Header.Get("Access-Token")
		json.NewDecoder(r.Body).Decode(&got)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	b := NewPushbulletBackend("test-key")
	b.url = ts.URL

	ev := alert.Event{Type: alert.NewPort, Port: 9090}
	if err := b.Send(ev); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotToken != "test-key" {
		t.Errorf("expected Access-Token test-key, got %s", gotToken)
	}
	if got["type"] != "note" {
		t.Errorf("expected type note, got %v", got["type"])
	}
}

func TestPushbulletBackendNon2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer ts.Close()

	b := NewPushbulletBackend("bad-key")
	b.url = ts.URL

	if err := b.Send(alert.Event{Type: alert.NewPort, Port: 80}); err == nil {
		t.Fatal("expected error for non-2xx status")
	}
}

func TestPushbulletBackendBadURL(t *testing.T) {
	b := NewPushbulletBackend("key")
	b.url = "http://127.0.0.1:0"

	if err := b.Send(alert.Event{Type: alert.NewPort, Port: 443}); err == nil {
		t.Fatal("expected error for unreachable URL")
	}
}
