package notify

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewRelicBackendName(t *testing.T) {
	b := NewNewRelicBackend("123", "key")
	if b.Name() != "newrelic" {
		t.Fatalf("expected newrelic, got %s", b.Name())
	}
}

func TestNewRelicBackendSendsJSON(t *testing.T) {
	var got []map[string]interface{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &got)
		if r.Header.Get("X-Insert-Key") != "insertkey" {
			t.Error("missing insert key header")
		}
		w.WriteHeader(200)
	}))
	defer ts.Close()
	b := &NewRelicBackend{accountID: "123", insertKey: "insertkey", url: ts.URL}
	if err := b.Send(Event{Port: 9090, Type: "new"}); err != nil {
		t.Fatal(err)
	}
	if len(got) == 0 {
		t.Fatal("no events sent")
	}
	if got[0]["eventType"] != "PortwatchEvent" {
		t.Errorf("unexpected eventType: %v", got[0]["eventType"])
	}
}

func TestNewRelicBackendNon2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
	}))
	defer ts.Close()
	b := &NewRelicBackend{accountID: "1", insertKey: "k", url: ts.URL}
	if err := b.Send(Event{Port: 22, Type: "closed"}); err == nil {
		t.Fatal("expected error")
	}
}

func TestNewRelicBackendBadURL(t *testing.T) {
	b := &NewRelicBackend{accountID: "1", insertKey: "k", url: "://bad"}
	if err := b.Send(Event{Port: 80, Type: "new"}); err == nil {
		t.Fatal("expected error")
	}
}
