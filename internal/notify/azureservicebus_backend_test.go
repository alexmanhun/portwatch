package notify

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/alert"
)

func TestAzureServiceBusBackendName(t *testing.T) {
	b := NewAzureServiceBusBackend("http://example.com", "SharedAccessSignature sr=...")
	if b.Name() != "azureservicebus" {
		t.Fatalf("expected azureservicebus, got %s", b.Name())
	}
}

func TestAzureServiceBusBackendSendsJSON(t *testing.T) {
	var gotBody []byte
	var gotAuth string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		gotBody, _ = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusCreated)
	}))
	defer ts.Close()

	b := NewAzureServiceBusBackend(ts.URL, "sas-token-value")
	err := b.Send(alert.Event{Type: alert.NewPort, Port: 8080})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if gotAuth != "sas-token-value" {
		t.Errorf("expected Authorization header sas-token-value, got %q", gotAuth)
	}

	var payload map[string]interface{}
	if err := json.Unmarshal(gotBody, &payload); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if payload["source"] != "portwatch" {
		t.Errorf("expected source=portwatch, got %v", payload["source"])
	}
	if int(payload["port"].(float64)) != 8080 {
		t.Errorf("expected port=8080, got %v", payload["port"])
	}
}

func TestAzureServiceBusBackendNon2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer ts.Close()

	b := NewAzureServiceBusBackend(ts.URL, "bad-token")
	err := b.Send(alert.Event{Type: alert.NewPort, Port: 443})
	if err == nil {
		t.Fatal("expected error for non-2xx response")
	}
}

func TestAzureServiceBusBackendBadURL(t *testing.T) {
	b := NewAzureServiceBusBackend("http://127.0.0.1:0", "token")
	err := b.Send(alert.Event{Type: alert.ClosedPort, Port: 22})
	if err == nil {
		t.Fatal("expected error for unreachable URL")
	}
}
