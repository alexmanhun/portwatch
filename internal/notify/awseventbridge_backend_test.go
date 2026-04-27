package notify

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAWSEventBridgeBackendName(t *testing.T) {
	b := NewAWSEventBridgeBackend("http://localhost", "", "", "")
	if b.Name() != "awseventbridge" {
		t.Fatalf("expected awseventbridge, got %s", b.Name())
	}
}

func TestAWSEventBridgeBackendSendsJSON(t *testing.T) {
	var received map[string]interface{}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &received)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	b := NewAWSEventBridgeBackend(ts.URL, "portwatch.test", "PortEvent", "secret")
	err := b.Send(Event{Port: 8080, Type: "new", Message: "port 8080 opened"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	entries, ok := received["Entries"].([]interface{})
	if !ok || len(entries) == 0 {
		t.Fatal("expected Entries array in payload")
	}

	entry := entries[0].(map[string]interface{})
	if entry["Source"] != "portwatch.test" {
		t.Errorf("unexpected Source: %v", entry["Source"])
	}
	if entry["DetailType"] != "PortEvent" {
		t.Errorf("unexpected DetailType: %v", entry["DetailType"])
	}
}

func TestAWSEventBridgeBackendApiKeyHeader(t *testing.T) {
	var gotKey string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotKey = r.Header.Get("X-Api-Key")
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	b := NewAWSEventBridgeBackend(ts.URL, "", "", "mykey")
	_ = b.Send(Event{Port: 443, Type: "closed", Message: "port 443 closed"})

	if gotKey != "mykey" {
		t.Errorf("expected X-Api-Key mykey, got %q", gotKey)
	}
}

func TestAWSEventBridgeBackendNon2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	b := NewAWSEventBridgeBackend(ts.URL, "", "", "")
	err := b.Send(Event{Port: 22, Type: "new", Message: "port 22 opened"})
	if err == nil {
		t.Fatal("expected error for non-2xx response")
	}
}

func TestAWSEventBridgeBackendBadURL(t *testing.T) {
	b := NewAWSEventBridgeBackend("http://127.0.0.1:0", "", "", "")
	err := b.Send(Event{Port: 80, Type: "new", Message: "port 80 opened"})
	if err == nil {
		t.Fatal("expected error for unreachable URL")
	}
}
