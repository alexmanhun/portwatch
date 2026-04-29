package notify

import (
	"sync/atomic"
	"testing"
	"time"

	"portwatch/internal/alert"
)

type countingBackend struct {
	name  string
	calls atomic.Int64
}

func (c *countingBackend) Name() string { return c.name }
func (c *countingBackend) Send(_ alert.Event) error {
	c.calls.Add(1)
	return nil
}

func TestThrottleBackendName(t *testing.T) {
	inner := &countingBackend{name: "mock"}
	tb := NewThrottleBackend(inner, time.Second)
	if tb.Name() != "throttle(mock)" {
		t.Errorf("unexpected name: %s", tb.Name())
	}
}

func TestThrottleBackendAllowsFirstEvent(t *testing.T) {
	inner := &countingBackend{name: "mock"}
	tb := NewThrottleBackend(inner, time.Second)
	ev := alert.NewPortEvent(8080)

	if err := tb.Send(ev); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if inner.calls.Load() != 1 {
		t.Errorf("expected 1 call, got %d", inner.calls.Load())
	}
}

func TestThrottleBackendSuppressesDuplicate(t *testing.T) {
	inner := &countingBackend{name: "mock"}
	tb := NewThrottleBackend(inner, time.Hour)
	ev := alert.NewPortEvent(9000)

	_ = tb.Send(ev)
	_ = tb.Send(ev)
	_ = tb.Send(ev)

	if inner.calls.Load() != 1 {
		t.Errorf("expected 1 call after throttle, got %d", inner.calls.Load())
	}
}

func TestThrottleBackendAllowsAfterCooldown(t *testing.T) {
	inner := &countingBackend{name: "mock"}
	tb := NewThrottleBackend(inner, 20*time.Millisecond)
	ev := alert.NewPortEvent(7070)

	_ = tb.Send(ev)
	time.Sleep(40 * time.Millisecond)
	_ = tb.Send(ev)

	if inner.calls.Load() != 2 {
		t.Errorf("expected 2 calls after cooldown, got %d", inner.calls.Load())
	}
}

func TestThrottleBackendDistinctPortsIndependent(t *testing.T) {
	inner := &countingBackend{name: "mock"}
	tb := NewThrottleBackend(inner, time.Hour)

	_ = tb.Send(alert.NewPortEvent(1111))
	_ = tb.Send(alert.NewPortEvent(2222))
	_ = tb.Send(alert.NewPortEvent(1111)) // suppressed

	if inner.calls.Load() != 2 {
		t.Errorf("expected 2 calls for distinct ports, got %d", inner.calls.Load())
	}
}
