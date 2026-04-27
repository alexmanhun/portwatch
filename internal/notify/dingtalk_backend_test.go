package notify

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/alert"
)

func TestDingTalkBackendName(t *testing.T) {
	b := NewDingTalkBackend("http://example.com/webhook")
	if b.Name() != "dingtalk" {
		t.Errorf("expected 'dingtalk', got %q", b.Name())
	}
}

func TestDingTalkBackendSendsJSON(t *testing.T) {
	var received map[string]interface{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &received)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	b := NewDingTalkBackend(ts.URL)
	event := alert.Event{Type: alert.NewPortEvent, Port: 9090}
	if err := b.Send(event); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received["msgtype"] != "text" {
		t.Errorf("expected msgtype=text, got %v", received["msgtype"])
	}
	textBlock, ok := received["text"].(map[string]interface{})
	if !ok || textBlock["content"] == "" {
		t.Error("expected non-empty content in text block")
	}
}

func TestDingTalkBackendNon2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
	}))
	defer ts.Close()

	b := NewDingTalkBackend(ts.URL)
	event := alert.Event{Type: alert.ClosedPortEvent, Port: 80}
	if err := b.Send(event); err == nil {
		t.Error("expected error on non-2xx response")
	}
}

func TestDingTalkBackendBadURL(t *testing.T) {
	b := NewDingTalkBackend("http://127.0.0.1:0/bad")
	event := alert.Event{Type: alert.NewPortEvent, Port: 3306}
	if err := b.Send(event); err == nil {
		t.Error("expected error for unreachable URL")
	}
}
