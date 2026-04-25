package notify

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/alert"
)

func TestWebexBackendName(t *testing.T) {
	b := NewWebexBackend("tok", "room1")
	if b.Name() != "webex" {
		t.Fatalf("expected 'webex', got %q", b.Name())
	}
}

func TestWebexBackendSendsJSON(t *testing.T) {
	var gotAuth, gotContentType string
	var gotBody map[string]string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		gotContentType = r.Header.Get("Content-Type")
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &gotBody)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	b := NewWebexBackend("mytoken", "roomABC")
	b.url = ts.URL

	ev := alert.Event{Type: alert.EventNewPort, Port: 8080, Message: "new port detected"}
	if err := b.Send(ev); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.HasPrefix(gotAuth, "Bearer ") {
		t.Errorf("expected Bearer auth, got %q", gotAuth)
	}
	if !strings.Contains(gotContentType, "application/json") {
		t.Errorf("expected JSON content type, got %q", gotContentType)
	}
	if gotBody["roomId"] != "roomABC" {
		t.Errorf("expected roomId 'roomABC', got %q", gotBody["roomId"])
	}
	if !strings.Contains(gotBody["text"], "8080") {
		t.Errorf("expected port 8080 in text, got %q", gotBody["text"])
	}
}

func TestWebexBackendNon2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer ts.Close()

	b := NewWebexBackend("bad-token", "room")
	b.url = ts.URL

	ev := alert.Event{Type: alert.EventNewPort, Port: 443, Message: "test"}
	if err := b.Send(ev); err == nil {
		t.Fatal("expected error for non-2xx response")
	}
}

func TestWebexBackendBadURL(t *testing.T) {
	b := NewWebexBackend("tok", "room")
	b.url = "://bad-url"

	ev := alert.Event{Type: alert.EventClosedPort, Port: 22, Message: "closed"}
	if err := b.Send(ev); err == nil {
		t.Fatal("expected error for bad URL")
	}
}
