package notify

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDingTalkBackendName(t *testing.T) {
	b := NewDingTalkBackend("http://example.com")
	if b.Name() != "dingtalk" {
		t.Fatalf("expected dingtalk, got %s", b.Name())
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
	err := b.Send(Event{Type: "new", Port: 9090})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received["msgtype"] != "text" {
		t.Errorf("expected msgtype=text, got %v", received["msgtype"])
	}
	textMap, ok := received["text"].(map[string]interface{})
	if !ok {
		t.Fatal("text field missing")
	}
	if textMap["content"] == "" {
		t.Error("content should not be empty")
	}
}

func TestDingTalkBackendNon2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	b := NewDingTalkBackend(ts.URL)
	if err := b.Send(Event{Type: "closed", Port: 80}); err == nil {
		t.Fatal("expected error for non-2xx response")
	}
}

func TestDingTalkBackendBadURL(t *testing.T) {
	b := NewDingTalkBackend("http://127.0.0.1:0")
	if err := b.Send(Event{Type: "new", Port: 443}); err == nil {
		t.Fatal("expected error for bad URL")
	}
}
