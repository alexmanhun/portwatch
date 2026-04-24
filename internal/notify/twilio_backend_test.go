package notify

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/alert"
)

func TestTwilioBackendName(t *testing.T) {
	b := NewTwilioBackend("sid", "token", "+10000000000", "+19999999999")
	if b.Name() != "twilio" {
		t.Fatalf("expected twilio, got %s", b.Name())
	}
}

func TestTwilioBackendSendsForm(t *testing.T) {
	var gotBody string
	var gotAuth string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw, _ := io.ReadAll(r.Body)
		gotBody = string(raw)
		gotAuth = r.Header.Get("Authorization")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"sid": "SM123"})
	}))
	defer ts.Close()

	b := NewTwilioBackend("ACtest", "authtoken", "+10000000000", "+19999999999")
	// Override the endpoint by pointing client at test server via transport trick.
	// We patch the URL directly on the backend for testing.
	b.client = ts.Client()
	// Re-route requests to test server.
	b.client.Transport = roundTripFunc(func(req *http.Request) (*http.Response, error) {
		req.URL, _ = url.Parse(ts.URL + req.URL.Path)
		return http.DefaultTransport.RoundTrip(req)
	})

	ev := alert.Event{Type: alert.EventNewPort, Port: 8080}
	if err := b.Send(ev); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(gotBody, "Body=") {
		t.Errorf("expected form body with Body= field, got: %s", gotBody)
	}
	if !strings.Contains(gotBody, "8080") {
		t.Errorf("expected port 8080 in body, got: %s", gotBody)
	}
	if gotAuth == "" {
		t.Error("expected Authorization header to be set")
	}
}

func TestTwilioBackendNon2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	b := NewTwilioBackend("ACtest", "bad", "+1", "+2")
	b.client = ts.Client()
	b.client.Transport = roundTripFunc(func(req *http.Request) (*http.Response, error) {
		req.URL, _ = url.Parse(ts.URL + req.URL.Path)
		return http.DefaultTransport.RoundTrip(req)
	})

	ev := alert.Event{Type: alert.EventNewPort, Port: 443}
	if err := b.Send(ev); err == nil {
		t.Fatal("expected error on non-2xx response")
	}
}

func TestTwilioBackendBadURL(t *testing.T) {
	b := NewTwilioBackend("AC bad sid", "token", "+1", "+2")
	b.client = &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			return nil, &url.Error{Op: "Post", URL: req.URL.String(), Err: io.EOF}
		}),
	}
	ev := alert.Event{Type: alert.EventClosedPort, Port: 22}
	if err := b.Send(ev); err == nil {
		t.Fatal("expected error on transport failure")
	}
}
