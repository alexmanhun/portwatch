package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/alert"
)

func TestFlowdockBackendName(t *testing.T) {
	b := NewFlowdockBackend("token123")
	if b.Name() != "flowdock" {
		t.Errorf("expected name 'flowdock', got %q", b.Name())
	}
}

func TestFlowdockBackendSendsJSON(t *testing.T) {
	var received map[string]interface{}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("failed to decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	// Override the URL by injecting a custom client that redirects to the test server.
	b := NewFlowdockBackend("mytoken")
	b.client = &http.Client{
		Transport: rewriteTransport(ts.URL),
	}

	event := alert.Event{Type: alert.EventNewPort, Port: 8080, Message: "port 8080 opened"}
	if err := b.Send(event); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	content, ok := received["content"].(string)
	if !ok || !strings.Contains(content, "port 8080 opened") {
		t.Errorf("expected content to contain event message, got %q", content)
	}

	if received["external_user_name"] != "portwatch" {
		t.Errorf("expected external_user_name 'portwatch', got %v", received["external_user_name"])
	}
}

func TestFlowdockBackendNon2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer ts.Close()

	b := NewFlowdockBackend("badtoken")
	b.client = &http.Client{Transport: rewriteTransport(ts.URL)}

	event := alert.Event{Type: alert.EventNewPort, Port: 9090, Message: "port 9090 opened"}
	err := b.Send(event)
	if err == nil {
		t.Fatal("expected error for non-2xx response")
	}
}

func TestFlowdockBackendBadURL(t *testing.T) {
	b := NewFlowdockBackend("")
	// Force a dial error by pointing to an invalid address.
	b.client = &http.Client{Transport: rewriteTransport("http://127.0.0.1:1")}

	event := alert.Event{Type: alert.EventClosedPort, Port: 22, Message: "port 22 closed"}
	err := b.Send(event)
	if err == nil {
		t.Fatal("expected error for unreachable URL")
	}
}
