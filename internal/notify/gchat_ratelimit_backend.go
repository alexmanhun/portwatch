package notify

import (
	"fmt"
	"sync"
	"time"

	"github.com/user/portwatch/internal/alert"
)

// RateLimitedBackend wraps any Backend and throttles Send calls.
// If more than MaxEvents events are dispatched within Window, extras are dropped.
type RateLimitedBackend struct {
	mu        sync.Mutex
	inner     Backend
	window    time.Duration
	maxEvents int
	counts    []time.Time
}

// NewRateLimitedBackend wraps inner, allowing at most maxEvents per window.
func NewRateLimitedBackend(inner Backend, maxEvents int, window time.Duration) *RateLimitedBackend {
	return &RateLimitedBackend{
		inner:     inner,
		maxEvents: maxEvents,
		window:    window,
	}
}

func (r *RateLimitedBackend) Name() string {
	return fmt.Sprintf("ratelimited(%s)", r.inner.Name())
}

func (r *RateLimitedBackend) Send(event alert.Event) error {
	r.mu.Lock()
	now := time.Now()
	cutoff := now.Add(-r.window)
	filtered := r.counts[:0]
	for _, t := range r.counts {
		if t.After(cutoff) {
			filtered = append(filtered, t)
		}
	}
	r.counts = filtered
	if len(r.counts) >= r.maxEvents {
		r.mu.Unlock()
		return fmt.Errorf("ratelimit: dropping event for backend %s", r.inner.Name())
	}
	r.counts = append(r.counts, now)
	r.mu.Unlock()
	return r.inner.Send(event)
}
