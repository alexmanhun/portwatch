package notify

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/alert"
)

func TestStatusPageBackendName(t *testing.T) {
	b := NewStatusPageBackend("key", "page", "comp", "")
	if b.Name() != "statuspage" {
		t.Fatalf("expected statuspage, got %s", b.Name())
	}
}

func TestStatusPageBackendSendsJSON(t *testing.T) {
	var gotBody map[string]interface{}
	var gotAuth string

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &gotBody)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	b := NewStatusPageBackend("mykey", "mypageid", "mycompid", srv.URL)
	e := alert.Event{Type: alert.EventNewPort, Port: 8080}

	if err := b.Send(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if gotAuth != "OAuth mykey" {
		t.Errorf("expected OAuth mykey, got %s", gotAuth)
	}

	comp, ok := gotBody["component"].(map[string]interface{})
	if !ok {
		t.Fatal("missing component key in payload")
	}
	if comp["status"] != "under_maintenance" {
		t.Errorf("expected under_maintenance, got %v", comp["status"])
	}
}

func TestStatusPageBackendClosedPortStatus(t *testing.T) {
	var gotBody map[string]interface{}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &gotBody)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	b := NewStatusPageBackend("k", "p", "c", srv.URL)
	e := alert.Event{Type: alert.EventClosedPort, Port: 443}

	if err := b.Send(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	comp := gotBody["component"].(map[string]interface{})
	if comp["status"] != "operational" {
		t.Errorf("expected operational, got %v", comp["status"])
	}
}

func TestStatusPageBackendNon2xx(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer srv.Close()

	b := NewStatusPageBackend("bad", "p", "c", srv.URL)
	if err := b.Send(alert.Event{Type: alert.EventNewPort}); err == nil {
		t.Fatal("expected error for non-2xx response")
	}
}

func TestStatusPageBackendBadURL(t *testing.T) {
	b := NewStatusPageBackend("k", "p", "c", "http://127.0.0.1:0")
	if err := b.Send(alert.Event{Type: alert.EventNewPort}); err == nil {
		t.Fatal("expected error for unreachable URL")
	}
}
