package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/alert"
)

func TestAWSSNSBackendName(t *testing.T) {
	b := NewAWSSNSBackend("arn:aws:sns:us-east-1:123456789012:portwatch", "us-east-1", "key", "secret")
	if b.Name() != "awssns" {
		t.Errorf("expected name 'awssns', got %q", b.Name())
	}
}

func TestAWSSNSBackendSendsJSON(t *testing.T) {
	var received map[string]interface{}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	b := NewAWSSNSBackend("arn:aws:sns:us-east-1:123456789012:portwatch", "us-east-1", "key", "secret")
	b.(*awsSNSBackend).endpoint = ts.URL

	ev := alert.NewPortEvent(8080)
	if err := b.Send(ev); err != nil {
		t.Fatalf("Send returned error: %v", err)
	}

	if received["TopicArn"] != "arn:aws:sns:us-east-1:123456789012:portwatch" {
		t.Errorf("unexpected TopicArn: %v", received["TopicArn"])
	}
	if received["Message"] == nil {
		t.Error("expected Message field in payload")
	}
}

func TestAWSSNSBackendNon2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	b := NewAWSSNSBackend("arn:aws:sns:us-east-1:123456789012:portwatch", "us-east-1", "key", "secret")
	b.(*awsSNSBackend).endpoint = ts.URL

	ev := alert.NewPortEvent(9090)
	if err := b.Send(ev); err == nil {
		t.Error("expected error on non-2xx response")
	}
}

func TestAWSSNSBackendBadURL(t *testing.T) {
	b := NewAWSSNSBackend("arn:aws:sns:us-east-1:123456789012:portwatch", "us-east-1", "key", "secret")
	b.(*awsSNSBackend).endpoint = "://invalid"

	ev := alert.NewPortEvent(443)
	if err := b.Send(ev); err == nil {
		t.Error("expected error on bad URL")
	}
}
