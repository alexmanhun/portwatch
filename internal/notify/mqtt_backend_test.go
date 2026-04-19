package notify

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMQTTBackendName(t *testing.T) {
	b := NewMQTTBackend("http://localhost:8080/publish", "portwatch/alerts")
	if b.Name() != "mqtt" {
		t.Fatalf("expected mqtt, got %s", b.Name())
	}
}

func TestMQTTBackendSendsJSON(t *testing.T) {
	var got map[string]interface{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &got)
		w.WriteHeader(200)
	}))
	defer ts.Close()

	b := NewMQTTBackend(ts.URL, "portwatch/alerts")
	err := b.Send(Event{Type: "new", Port: 8080})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["topic"] != "portwatch/alerts" {
		t.Errorf("expected topic portwatch/alerts, got %v", got["topic"])
	}
}

func TestMQTTBackendNon2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(503)
	}))
	defer ts.Close()

	b := NewMQTTBackend(ts.URL, "alerts")
	if err := b.Send(Event{Type: "closed", Port: 22}); err == nil {
		t.Fatal("expected error on non-2xx")
	}
}

func TestMQTTBackendBadURL(t *testing.T) {
	b := NewMQTTBackend("http://127.0.0.1:0/publish", "alerts")
	if err := b.Send(Event{Type: "new", Port: 9090}); err == nil {
		t.Fatal("expected error on bad URL")
	}
}
