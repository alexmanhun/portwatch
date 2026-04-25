package notify

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestClickSendBackendName(t *testing.T) {
	b := NewClickSendBackend("user", "key", "+15550001234")
	if b.Name() != "clicksend" {
		t.Fatalf("expected clicksend, got %s", b.Name())
	}
}

func TestClickSendBackendSendsJSON(t *testing.T) {
	var captured map[string]interface{}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		ct := r.Header.Get("Content-Type")
		if !strings.Contains(ct, "application/json") {
			t.Errorf("unexpected Content-Type: %s", ct)
		}
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &captured)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	b := NewClickSendBackend("user", "key", "+15550001234")
	// Override the API URL by pointing the client at our test server via a
	// custom transport that rewrites the host.
	b.client = ts.Client()
	// Since we cannot easily override the hardcoded URL in a unit test without
	// refactoring, we instead verify the struct is constructed correctly and
	// that a real HTTP call with a mock server returns no error when status 200.
	// We replace the internal client with one using a RoundTripper redirect.
	b.client = &http.Client{
		Transport: &rewriteTransport{target: ts.URL},
	}

	evt := Event{Type: "new_port", Port: 8080}
	if err := b.Send(evt); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestClickSendBackendNon2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	b := NewClickSendBackend("user", "badkey", "+15550001234")
	b.client = &http.Client{
		Transport: &rewriteTransport{target: ts.URL},
	}

	if err := b.Send(Event{Type: "new_port", Port: 443}); err == nil {
		t.Fatal("expected error for non-2xx response")
	}
}

func TestClickSendBackendBadURL(t *testing.T) {
	b := NewClickSendBackend("user", "key", "+15550001234")
	b.client = &http.Client{
		Transport: &rewriteTransport{target: "http://127.0.0.1:1"},
	}
	if err := b.Send(Event{Type: "closed_port", Port: 22}); err == nil {
		t.Fatal("expected connection error")
	}
}

// rewriteTransport redirects all requests to a fixed target host for testing.
type rewriteTransport struct{ target string }

func (rt *rewriteTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req2 := req.Clone(req.Context())
	target := rt.target
	req2.URL.Scheme = "http"
	parsed, _ := http.NewRequest("", target, nil)
	req2.URL.Host = parsed.URL.Host
	return http.DefaultTransport.RoundTrip(req2)
}
