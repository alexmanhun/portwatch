package notify

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestXMPPBackendName(t *testing.T) {
	b := NewXMPPBackend("http://localhost", "user@example.com")
	if b.Name() != "xmpp" {
		t.Fatalf("expected xmpp, got %s", b.Name())
	}
}

func TestXMPPBackendSendsForm(t *testing.T) {
	var gotTo, gotBody string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		gotTo = r.FormValue("to")
		gotBody = r.FormValue("body")
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	b := NewXMPPBackend(ts.URL, "alice@example.com")
	err := b.Send(Event{Type: "new", Port: 8080})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotTo != "alice@example.com" {
		t.Errorf("expected to=alice@example.com, got %s", gotTo)
	}
	if gotBody == "" {
		t.Error("expected non-empty body")
	}
}

func TestXMPPBackendNon2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer ts.Close()

	b := NewXMPPBackend(ts.URL, "bob@example.com")
	if err := b.Send(Event{Type: "closed", Port: 22}); err == nil {
		t.Fatal("expected error for non-2xx response")
	}
}

func TestXMPPBackendBadURL(t *testing.T) {
	b := NewXMPPBackend("http://127.0.0.1:0", "x@y.com")
	if err := b.Send(Event{Type: "new", Port: 9000}); err == nil {
		t.Fatal("expected error for bad URL")
	}
}
