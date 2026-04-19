package notify

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestKafkaBackendName(t *testing.T) {
	b := NewKafkaBackend("http://localhost:8082", "portwatch")
	if b.Name() != "kafka" {
		t.Fatalf("expected kafka, got %s", b.Name())
	}
}

func TestKafkaBackendSendsJSON(t *testing.T) {
	var got map[string]any
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &got)
		if ct := r.Header.Get("Content-Type"); ct != "application/vnd.kafka.json.v2+json" {
			t.Errorf("unexpected content-type: %s", ct)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	b := NewKafkaBackend(ts.URL, "portwatch")
	if err := b.Send(Event{Type: "new_port", Port: 8080}); err != nil {
		t.Fatal(err)
	}
	records, ok := got["records"].([]any)
	if !ok || len(records) == 0 {
		t.Fatal("expected records in payload")
	}
}

func TestKafkaBackendNon2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer ts.Close()

	b := NewKafkaBackend(ts.URL, "portwatch")
	if err := b.Send(Event{Type: "new_port", Port: 9090}); err == nil {
		t.Fatal("expected error on non-2xx")
	}
}

func TestKafkaBackendBadURL(t *testing.T) {
	b := NewKafkaBackend("http://127.0.0.1:1", "portwatch")
	if err := b.Send(Event{Type: "new_port", Port: 22}); err == nil {
		t.Fatal("expected error on bad URL")
	}
}
