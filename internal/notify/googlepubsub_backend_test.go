package notify

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
)

func TestGooglePubSubBackendName(t *testing.T) {
	b := NewGooglePubSubBackend("proj", "topic", "tok")
	if b.Name() != "googlepubsub" {
		t.Fatalf("expected googlepubsub, got %s", b.Name())
	}
}

func TestGooglePubSubBackendSendsJSON(t *testing.T) {
	var captured map[string]interface{}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer mytoken" {
			t.Errorf("missing or wrong Authorization header")
		}
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &captured)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{}`))
	}))
	defer ts.Close()

	b := NewGooglePubSubBackend("my-project", "my-topic", "mytoken")
	// Override URL by pointing to test server via a custom transport isn't
	// straightforward, so we swap the client with one that redirects to ts.
	b.client = ts.Client()
	// Patch the URL by overriding projectID/topicID to craft a path the test
	// server will accept (the server ignores path).
	b.projectID = ""
	b.topicID = ""

	// Use a real server URL directly.
	b2 := &GooglePubSubBackend{
		projectID: "proj",
		topicID:   "topic",
		token:     "mytoken",
		client:    &http.Client{Transport: &redirectTransport{base: ts.URL}},
	}

	event := alert.Event{Port: 8080, Type: alert.NewPort, Timestamp: time.Now()}
	if err := b2.Send(event); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGooglePubSubBackendNon2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	b := &GooglePubSubBackend{
		projectID: "proj",
		topicID:   "topic",
		token:     "tok",
		client:    &http.Client{Transport: &redirectTransport{base: ts.URL}},
	}
	event := alert.Event{Port: 443, Type: alert.ClosedPort, Timestamp: time.Now()}
	if err := b.Send(event); err == nil {
		t.Fatal("expected error for non-2xx response")
	}
}

func TestGooglePubSubBackendBadURL(t *testing.T) {
	b := &GooglePubSubBackend{
		projectID: "p",
		topicID:   "t",
		token:     "tok",
		client:    &http.Client{},
	}
	// Force a dial error by using an invalid host.
	b.projectID = "\x00invalid"
	event := alert.Event{Port: 22, Type: alert.NewPort, Timestamp: time.Now()}
	if err := b.Send(event); err == nil {
		t.Fatal("expected error for bad URL")
	}
}

// redirectTransport rewrites every request to target the given base URL,
// preserving the path. This lets tests redirect production URLs to httptest.
type redirectTransport struct {
	base string
}

func (rt *redirectTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req2 := req.Clone(req.Context())
	req2.URL.Scheme = "http"
	req2.URL.Host = rt.base[len("http://"):]
	return http.DefaultTransport.RoundTrip(req2)
}
