package notify

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

// stubBackend records sent messages.
type stubBackend struct {
	name string
	msgs []Message
	err  error
}

func (s *stubBackend) Name() string { return s.name }
func (s *stubBackend) Send(m Message) error {
	s.msgs = append(s.msgs, m)
	return s.err
}

func TestDispatchSendsToAllBackends(t *testing.T) {
	a := &stubBackend{name: "a"}
	b := &stubBackend{name: "b"}
	d := New(a, b)
	msg := Message{Level: LevelAlert, Title: "test", Body: "port opened", Port: 8080, Event: "new"}
	d.Dispatch(msg)
	if len(a.msgs) != 1 || len(b.msgs) != 1 {
		t.Fatalf("expected 1 message each, got a=%d b=%d", len(a.msgs), len(b.msgs))
	}
	if a.msgs[0].Port != 8080 {
		t.Errorf("unexpected port: %d", a.msgs[0].Port)
	}
}

func TestDispatchLogsBackendError(t *testing.T) {
	bad := &stubBackend{name: "bad", err: errors.New("boom")}
	var errBuf bytes.Buffer
	d := New(bad)
	d.SetErrOut(&errBuf)
	d.Dispatch(Message{Level: LevelWarn, Title: "x", Port: 22})
	if !strings.Contains(errBuf.String(), "bad") {
		t.Errorf("expected backend name in error output, got: %s", errBuf.String())
	}
}

func TestLogBackendOutput(t *testing.T) {
	var buf bytes.Buffer
	lb := NewLogBackend(&buf)
	err := lb.Send(Message{Level: LevelInfo, Title: "port closed", Body: "gone", Port: 443})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "443") || !strings.Contains(out, "INFO") {
		t.Errorf("unexpected log output: %s", out)
	}
}
