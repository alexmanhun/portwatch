package notify

import (
	"log/syslog"
	"testing"
)

func TestSyslogBackendName(t *testing.T) {
	b, err := NewSyslogBackend("portwatch-test", syslog.LOG_INFO|syslog.LOG_DAEMON)
	if err != nil {
		t.Skip("syslog not available:", err)
	}
	defer b.Close()
	if b.Name() != "syslog" {
		t.Errorf("expected syslog, got %s", b.Name())
	}
}

func TestSyslogBackendSendsEvent(t *testing.T) {
	b, err := NewSyslogBackend("portwatch-test", syslog.LOG_INFO|syslog.LOG_DAEMON)
	if err != nil {
		t.Skip("syslog not available:", err)
	}
	defer b.Close()

	event := Event{Kind: "new", Port: 8080, Message: "port opened"}
	if err := b.Send(event); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestSyslogBackendCloseNil(t *testing.T) {
	b := &SyslogBackend{}
	if err := b.Close(); err != nil {
		t.Errorf("expected nil error on close with nil writer, got %v", err)
	}
}

func TestSyslogBackendInvalidTag(t *testing.T) {
	// This simply exercises the constructor path; skip if syslog unavailable.
	_, err := NewSyslogBackend("", syslog.LOG_DEBUG|syslog.LOG_LOCAL0)
	if err != nil {
		t.Skip("syslog not available:", err)
	}
}
