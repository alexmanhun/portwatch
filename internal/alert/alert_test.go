package alert_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"portwatch/internal/alert"
)

func TestNotifyWritesOutput(t *testing.T) {
	var buf bytes.Buffer
	n := alert.New(&buf)

	e := alert.Event{
		Timestamp: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		Level:     alert.LevelAlert,
		Port:      8080,
		Message:   "new open port detected: 8080",
	}

	if err := n.Notify(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "ALERT") {
		t.Errorf("expected ALERT in output, got: %s", out)
	}
	if !strings.Contains(out, "8080") {
		t.Errorf("expected port 8080 in output, got: %s", out)
	}
}

func TestNewPortEvent(t *testing.T) {
	e := alert.NewPortEvent(9090)
	if e.Level != alert.LevelAlert {
		t.Errorf("expected ALERT level, got %s", e.Level)
	}
	if e.Port != 9090 {
		t.Errorf("expected port 9090, got %d", e.Port)
	}
	if !strings.Contains(e.Message, "9090") {
		t.Errorf("expected port in message, got: %s", e.Message)
	}
}

func TestClosedPortEvent(t *testing.T) {
	e := alert.ClosedPortEvent(443)
	if e.Level != alert.LevelWarn {
		t.Errorf("expected WARN level, got %s", e.Level)
	}
	if e.Port != 443 {
		t.Errorf("expected port 443, got %d", e.Port)
	}
}

func TestNotifyDefaultsToStdout(t *testing.T) {
	// Ensure New(nil) does not panic.
	n := alert.New(nil)
	if n == nil {
		t.Fatal("expected non-nil Notifier")
	}
}
