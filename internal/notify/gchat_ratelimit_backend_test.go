package notify

import (
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
)

type countingBackend struct {
	name  string
	calls atomic.Int32
}

func (c *countingBackend) Name() string { return c.name }
func (c *countingBackend) Send(_ alert.Event) error {
	c.calls.Add(1)
	return nil
}

func TestRateLimitedBackendName(t *testing.T) {
	inner := &countingBackend{name: "mock"}
	b := NewRateLimitedBackend(inner, 5, time.Minute)
	if b.Name() != "ratelimited(mock)" {
		t.Fatalf("unexpected name: %s", b.Name())
	}
}

func TestRateLimitedBackendAllowsUnderLimit(t *testing.T) {
	inner := &countingBackend{name: "mock"}
	b := NewRateLimitedBackend(inner, 3, time.Minute)
	ev := alert.Event{Type: alert.NewPort, Port: 80}
	for i := 0; i < 3; i++ {
		if err := b.Send(ev); err != nil {
			t.Fatalf("unexpected error on send %d: %v", i, err)
		}
	}
	if inner.calls.Load() != 3 {
		t.Fatalf("expected 3 calls, got %d", inner.calls.Load())
	}
}

func TestRateLimitedBackendDropsOverLimit(t *testing.T) {
	inner := &countingBackend{name: "mock"}
	b := NewRateLimitedBackend(inner, 2, time.Minute)
	ev := alert.Event{Type: alert.NewPort, Port: 443}
	b.Send(ev)
	b.Send(ev)
	err := b.Send(ev)
	if err == nil {
		t.Fatal("expected rate limit error on third send")
	}
	if inner.calls.Load() != 2 {
		t.Fatalf("expected 2 calls, got %d", inner.calls.Load())
	}
}

func TestRateLimitedBackendResetsAfterWindow(t *testing.T) {
	inner := &countingBackend{name: "mock"}
	b := NewRateLimitedBackend(inner, 1, 50*time.Millisecond)
	ev := alert.Event{Type: alert.NewPort, Port: 22}
	if err := b.Send(ev); err != nil {
		t.Fatalf("first send failed: %v", err)
	}
	if err := b.Send(ev); !errors.Is(err, err) || err == nil {
		t.Fatal("expected rate limit on second send")
	}
	time.Sleep(60 * time.Millisecond)
	if err := b.Send(ev); err != nil {
		t.Fatalf("expected send to succeed after window reset: %v", err)
	}
}
