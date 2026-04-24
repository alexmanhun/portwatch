package notify

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/alert"
)

func TestLarkBackendName(t *testing.T) {
	b := NewLarkBackend("https://example.com/hook")
	if b.Name() != "lark" {
		t.Fatalf("expected name 'lark', got %q", b.Name())
	}
}

func TestLarkBackendSendsJSON(t *testing.T) {
	var received map[string]interface{}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &received)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	b := NewLarkBackend(ts.URL)
	event := alert.Event{Type: alert.EventNewPort, Port: 9090, Message: "new port detected"}

	if err := b.Send(event); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if received["msg_type"] != "text" {
		t.Errorf("expected msg_type 'text', got %v", received["msg_type"])
	}

	content, ok := received["content"].(map[string]interface{})
	if !ok {
		t.Fatal("expected content to be a map")
	}
	text, _ := content["text"].(string)
	if text == "" {
		t.Error("expected non-empty text in content")
	}
}

func TestLarkBackendNon2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	b := NewLarkBackend(ts.URL)
	event := alert.Event{Type: alert.EventNewPort, Port: 8080, Message: "test"}

	if err := b.Send(event); err == nil {
		t.Fatal("expected error for non-2xx response")
	}
}

func TestLarkBackendBadURL(t *testing.T) {
	b := NewLarkBackend("://bad-url")
	event := alert.Event{Type: alert.EventClosedPort, Port: 443, Message: "closed"}

	if err := b.Send(event); err == nil {
		t.Fatal("expected error for invalid URL")
	}
}
