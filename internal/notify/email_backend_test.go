package notify

import (
	"fmt"
	"net/smtp"
	"strings"
	"testing"
)

func TestEmailBackendName(t *testing.T) {
	b := NewEmailBackend(EmailConfig{})
	if b.Name() != "email" {
		t.Errorf("expected 'email', got %q", b.Name())
	}
}

func TestEmailBackendSendsMessage(t *testing.T) {
	var capturedAddr string
	var capturedFrom string
	var capturedTo []string
	var capturedMsg string

	cfg := EmailConfig{
		Host:     "smtp.example.com",
		Port:     587,
		Username: "user",
		Password: "pass",
		From:     "portwatch@example.com",
		To:       []string{"admin@example.com"},
	}

	eb := &emailBackend{
		cfg: cfg,
		send: func(addr, from string, to []string, msg []byte) error {
			capturedAddr = addr
			capturedFrom = from
			capturedTo = to
			capturedMsg = string(msg)
			return nil
		},
	}

	if err := eb.Send("new_port", "port 8080 opened"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if capturedAddr != "smtp.example.com:587" {
		t.Errorf("wrong addr: %s", capturedAddr)
	}
	if capturedFrom != cfg.From {
		t.Errorf("wrong from: %s", capturedFrom)
	}
	if len(capturedTo) != 1 || capturedTo[0] != "admin@example.com" {
		t.Errorf("wrong to: %v", capturedTo)
	}
	if !strings.Contains(capturedMsg, "new_port") {
		t.Error("message missing event")
	}
	if !strings.Contains(capturedMsg, "port 8080 opened") {
		t.Error("message missing detail")
	}
}

func TestEmailBackendSendError(t *testing.T) {
	eb := &emailBackend{
		cfg: EmailConfig{Host: "bad", Port: 25, To: []string{"x@x.com"}},
		send: func(addr, from string, to []string, msg []byte) error {
			return fmt.Errorf("connection refused")
		},
	}
	err := eb.Send("closed_port", "port 22 closed")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	_ = smtp.PlainAuth // ensure smtp imported
}
