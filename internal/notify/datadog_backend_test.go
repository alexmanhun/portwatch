package notify

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDatadogBackendName(t *testing.T) {
	b := NewDatadogBackend("key")
	if b.Name() != "datadog" {
		t.Fatalf("expected datadog, got %s", b.Name())
	}
}

func TestDatadogBackendSendsJSON(t *testing.T) {
	var got map[string]interface{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &got)
		if r.Header.Get("DD-API-KEY") != "testkey" {
			t.Error("missing api key header")
		}
		w.WriteHeader(202)
	}))
	defer ts.Close()
	b := &DatadogBackend{apiKey: "testkey", url: ts.URL}
	if err := b.Send(Event{Port: 8080, Type: "new"}); err != nil {
		t.Fatal(err)
	}
	if got["title"] != "portwatch: new" {
		t.Errorf("unexpected title: %v", got["title"])
	}
}

func TestDatadogBackendNon2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(403)
	}))
	defer ts.Close()
	b := &DatadogBackend{apiKey: "k", url: ts.URL}
	if err := b.Send(Event{Port: 22, Type: "closed"}); err == nil {
		t.Fatal("expected error")
	}
}

func TestDatadogBackendBadURL(t *testing.T) {
	b := &DatadogBackend{apiKey: "k", url: "://bad"}
	if err := b.Send(Event{Port: 80, Type: "new"}); err == nil {
		t.Fatal("expected error")
	}
}
