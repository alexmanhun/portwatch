package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/alert"
)

func TestSplunkBackendName(t *testing.T) {
	b := NewSplunkBackend("http://localhost", "token")
	if b.Name() != "splunk" {
		t.Fatalf("expected splunk, got %s", b.Name())
	}
}

func TestSplunkBackendSendsJSON(t *testing.T) {
	var got map[string]interface{}
	var authHeader string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader = r.Header.Get("Authorization")
		json.NewDecoder(r.Body).Decode(&got)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	b := NewSplunkBackend(ts.URL, "mytoken")
	e := alert.Event{Port: 8080, Type: "new", Message: "port opened"}
	if err := b.Send(e); err != nil {
		t.Fatal(err)
	}
	if authHeader != "Splunk mytoken" {
		t.Errorf("unexpected auth header: %s", authHeader)
	}
	ev, ok := got["event"].(map[string]interface{})
	if !ok {
		t.Fatal("missing event field")
	}
	if int(ev["port"].(float64)) != 8080 {
		t.Errorf("unexpected port: %v", ev["port"])
	}
	if got["sourcetype"] != "portwatch" {
		t.Errorf("unexpected sourcetype: %v", got["sourcetype"])
	}
}

func TestSplunkBackendNon2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	b := NewSplunkBackend(ts.URL, "tok")
	if err := b.Send(alert.Event{Port: 443, Type: "closed"}); err == nil {
		t.Fatal("expected error for non-2xx response")
	}
}

func TestSplunkBackendBadURL(t *testing.T) {
	b := NewSplunkBackend("://bad-url", "tok")
	if err := b.Send(alert.Event{Port: 80, Type: "new"}); err == nil {
		t.Fatal("expected error for bad URL")
	}
}
