package notify

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNtfyBackendName(t *testing.T) {
	n := NewNtfyBackend("https://ntfy.sh/test")
	if n.Name() != "ntfy" {
		t.Errorf("expected ntfy, got %s", n.Name())
	}
}

func TestNtfyBackendSendsMessage(t *testing.T) {
	var body string
	var title string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		body = string(b)
		title = r.Header.Get("Title")
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewNtfyBackend(ts.URL)
	err := n.Send(Event{Type: "closed_port", Message: "port 22 closed"})
	if err != nil {
		t.Fatal(err)
	}
	if body != "port 22 closed" {
		t.Errorf("unexpected body: %s", body)
	}
	if title != "portwatch: closed_port" {
		t.Errorf("unexpected title: %s", title)
	}
}

func TestNtfyBackendNon2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	n := NewNtfyBackend(ts.URL)
	if err := n.Send(Event{Type: "new_port", Message: "x"}); err == nil {
		t.Error("expected error")
	}
}

func TestNtfyBackendBadURL(t *testing.T) {
	n := NewNtfyBackend("http://127.0.0.1:0")
	if err := n.Send(Event{Type: "new_port", Message: "x"}); err == nil {
		t.Error("expected connection error")
	}
}
