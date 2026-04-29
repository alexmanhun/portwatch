package notify

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAWSLambdaBackendName(t *testing.T) {
	b := NewAWSLambdaBackend("https://example.lambda-url.us-east-1.on.aws/", "")
	if b.Name() != "awslambda" {
		t.Fatalf("expected awslambda, got %s", b.Name())
	}
}

func TestAWSLambdaBackendSendsJSON(t *testing.T) {
	var received map[string]interface{}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &received)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	b := NewAWSLambdaBackend(ts.URL, "")
	err := b.Send(Event{Type: "new_port", Port: 8080, Message: "port opened"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if received["source"] != "portwatch" {
		t.Errorf("expected source=portwatch, got %v", received["source"])
	}
	if received["event"] != "new_port" {
		t.Errorf("expected event=new_port, got %v", received["event"])
	}
	if int(received["port"].(float64)) != 8080 {
		t.Errorf("expected port=8080, got %v", received["port"])
	}
}

func TestAWSLambdaBackendApiKeyHeader(t *testing.T) {
	var gotKey string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotKey = r.Header.Get("x-api-key")
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	b := NewAWSLambdaBackend(ts.URL, "secret-key")
	_ = b.Send(Event{Type: "new_port", Port: 443, Message: "port opened"})

	if gotKey != "secret-key" {
		t.Errorf("expected x-api-key=secret-key, got %q", gotKey)
	}
}

func TestAWSLambdaBackendNon2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	b := NewAWSLambdaBackend(ts.URL, "")
	if err := b.Send(Event{Type: "new_port", Port: 80, Message: "port opened"}); err == nil {
		t.Fatal("expected error for non-2xx status")
	}
}

func TestAWSLambdaBackendBadURL(t *testing.T) {
	b := NewAWSLambdaBackend("http://127.0.0.1:0", "")
	if err := b.Send(Event{Type: "new_port", Port: 22, Message: "port opened"}); err == nil {
		t.Fatal("expected error for unreachable URL")
	}
}
