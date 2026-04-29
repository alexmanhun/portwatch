package notify

import (
	"fmt"
	"sync"
	"time"

	"github.com/user/portwatch/internal/alert"
)

// dedupKey uniquely identifies an event for deduplication purposes.
type dedupKey struct {
	port  int
	event string
}

// dedupEntry holds the last seen time for a dedup key.
type dedupEntry struct {
	seenAt time.Time
}

// DedupBackend wraps another Backend and suppresses duplicate events
// for the same port+event combination within a configurable window.
// Unlike ThrottleBackend (which suppresses per cooldown period),
// DedupBackend tracks the exact last-seen time and only forwards
// an event if the dedup window has expired since the last identical event.
type DedupBackend struct {
	inner  Backend
	window time.Duration
	mu     sync.Mutex
	seen   map[dedupKey]dedupEntry
}

// NewDedupBackend creates a DedupBackend wrapping inner.
// Events with the same port and event type seen within window duration
// are silently dropped; only the first occurrence is forwarded.
func NewDedupBackend(inner Backend, window time.Duration) *DedupBackend {
	return &DedupBackend{
		inner:  inner,
		window: window,
		seen:   make(map[dedupKey]dedupEntry),
	}
}

// Name returns a descriptive name including the wrapped backend's name.
func (d *DedupBackend) Name() string {
	return fmt.Sprintf("dedup(%s,window=%s)", d.inner.Name(), d.window)
}

// Send forwards the event to the inner backend only if no identical
// event (same port + event type) has been seen within the dedup window.
// Thread-safe.
func (d *DedupBackend) Send(e alert.Event) error {
	key := dedupKey{port: e.Port, event: e.Type}
	now := time.Now()

	d.mu.Lock()
	entry, exists := d.seen[key]
	if exists && now.Sub(entry.seenAt) < d.window {
		d.mu.Unlock()
		return nil // duplicate within window, drop silently
	}
	d.seen[key] = dedupEntry{seenAt: now}
	d.mu.Unlock()

	return d.inner.Send(e)
}

// Flush removes all dedup entries whose window has expired.
// Call periodically to prevent unbounded memory growth in long-running daemons.
func (d *DedupBackend) Flush() {
	now := time.Now()
	d.mu.Lock()
	defer d.mu.Unlock()
	for k, entry := range d.seen {
		if now.Sub(entry.seenAt) >= d.window {
			delete(d.seen, k)
		}
	}
}

// Close delegates to the inner backend if it implements io.Closer.
func (d *DedupBackend) Close() error {
	type closer interface {
		Close() error
	}
	if c, ok := d.inner.(closer); ok {
		return c.Close()
	}
	return nil
}
