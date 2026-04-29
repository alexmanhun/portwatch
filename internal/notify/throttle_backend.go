package notify

import (
	"fmt"
	"sync"
	"time"

	"portwatch/internal/alert"
)

// ThrottleBackend wraps another Backend and suppresses duplicate events
// for the same port within a configurable cooldown window.
type ThrottleBackend struct {
	inner    Backend
	cooldown time.Duration
	mu       sync.Mutex
	lastSent map[string]time.Time
}

// NewThrottleBackend creates a ThrottleBackend that delegates to inner
// but drops events for a given port+event key until cooldown has elapsed.
func NewThrottleBackend(inner Backend, cooldown time.Duration) *ThrottleBackend {
	return &ThrottleBackend{
		inner:    inner,
		cooldown: cooldown,
		lastSent: make(map[string]time.Time),
	}
}

func (t *ThrottleBackend) Name() string {
	return fmt.Sprintf("throttle(%s)", t.inner.Name())
}

func (t *ThrottleBackend) Send(ev alert.Event) error {
	key := fmt.Sprintf("%s:%d", ev.Type, ev.Port)

	t.mu.Lock()
	last, seen := t.lastSent[key]
	if seen && time.Since(last) < t.cooldown {
		t.mu.Unlock()
		return nil
	}
	t.lastSent[key] = time.Now()
	t.mu.Unlock()

	return t.inner.Send(ev)
}

func (t *ThrottleBackend) Close() error {
	if c, ok := t.inner.(interface{ Close() error }); ok {
		return c.Close()
	}
	return nil
}
