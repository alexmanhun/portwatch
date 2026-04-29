package notify

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/alert"
)

func TestClickSendBackendName(t *testing.T) {
	b := NewClickSendBackend("http://example.com", "user", "key", "+1555000000")
	if b.Name() != "clicksend" {
		t.Fatalf("expected clicksend, got %s", b.Name())
	}
}

func TestClickSendBackendSendsJSON(t *testing.T) {
	var got map[string]interface{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &got)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	b := NewClickSendBackend(ts.URL, "user", "key", "+1555000000")
	evt := alert.Event{Type: alert.EventNewPort, Port: 8080, Message: "new port"}
	if err := b.Send(evt); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	msgs, ok := got["messages"].([]interface{})
	if !ok || len(msgs) == 0 {
		t.Fatal("expected messages array")
	}
	msg := msgs[0].(map[string]interface{})
	if msg["to"] != "+1555000000" {
		t.Errorf("unexpected to: %v", msg["to"])
	}
}

func TestClickSendBackendNon2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer ts.Close()

	b := NewClickSendBackend(ts.URL, "user", "badkey", "+1555000000")
	evt := alert.Event{Type: alert.EventNewPort, Port: 9090, Message: "test"}
	if err := b.Send(evt); err == nil {
		t.Fatal("expected error for non-2xx response")
	}
}

func TestClickSendBackendBadURL(t *testing.T) {
	b := NewClickSendBackend("http://127.0.0.1:0", "user", "key", "+1555000000")
	evt := alert.Event{Type: alert.EventClosedPort, Port: 443, Message: "closed"}
	if err := b.Send(evt); err == nil {
		t.Fatal("expected error for unreachable URL")
	}
}
