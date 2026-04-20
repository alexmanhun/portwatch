package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGrafanaBackendName(t *testing.T) {
	b := NewGrafanaBackend("http://localhost:3000", "token")
	if b.Name() != "grafana" {
		t.Fatalf("expected grafana, got %s", b.Name())
	}
}

func TestGrafanaBackendSendsJSON(t *testing.T) {
	var received map[string]interface{}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/annotations" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer testkey" {
			t.Errorf("missing or wrong Authorization header")
		}
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	b := NewGrafanaBackend(ts.URL, "testkey")
	err := b.Send(Event{Type: "new_port", Port: 8080})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	text, _ := received["text"].(string)
	if text == "" {
		t.Error("expected non-empty text field")
	}
	tags, _ := received["tags"].([]interface{})
	if len(tags) == 0 {
		t.Error("expected tags to be populated")
	}
}

func TestGrafanaBackendNon2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer ts.Close()

	b := NewGrafanaBackend(ts.URL, "badkey")
	err := b.Send(Event{Type: "closed_port", Port: 443})
	if err == nil {
		t.Fatal("expected error for non-2xx response")
	}
}

func TestGrafanaBackendBadURL(t *testing.T) {
	b := NewGrafanaBackend("http://127.0.0.1:0", "key")
	err := b.Send(Event{Type: "new_port", Port: 22})
	if err == nil {
		t.Fatal("expected error for unreachable URL")
	}
}
