package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWebhookBackendSendsJSON(t *testing.T) {
	var received webhookPayload
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			http.Error(w, "bad json", 400)
			return
		}
		w.WriteHeader(200)
	}))
	defer ts.Close()

	wb := NewWebhookBackend(ts.URL)
	err := wb.Send(Message{Level: LevelAlert, Title: "new port", Body: "detected", Port: 9090, Event: "new"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.Port != 9090 {
		t.Errorf("expected port 9090, got %d", received.Port)
	}
	if received.Level != "ALERT" {
		t.Errorf("expected level ALERT, got %s", received.Level)
	}
}

func TestWebhookBackendNon2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
	}))
	defer ts.Close()

	wb := NewWebhookBackend(ts.URL)
	err := wb.Send(Message{Level: LevelWarn, Title: "x", Port: 80})
	if err == nil {
		t.Error("expected error for non-2xx response")
	}
}

func TestWebhookBackendName(t *testing.T) {
	wb := NewWebhookBackend("http://example.com")
	if wb.Name() != "webhook" {
		t.Errorf("unexpected name: %s", wb.Name())
	}
}
