package history

import "time"

// Filter holds optional criteria for querying history entries.
type Filter struct {
	Port  int       // 0 means any port
	Since time.Time // zero means no lower bound
	Event string    // "" means any event
}

// Query returns entries matching the given filter.
func (h *History) Query(f Filter) []Entry {
	h.mu.Lock()
	defer h.mu.Unlock()

	var result []Entry
	for _, e := range h.entries {
		if f.Port != 0 && e.Port != f.Port {
			continue
		}
		if !f.Since.IsZero() && e.Timestamp.Before(f.Since) {
			continue
		}
		if f.Event != "" && e.Event != f.Event {
			continue
		}
		result = append(result, e)
	}
	return result
}

// Last returns the most recent n entries, or all if n <= 0.
func (h *History) Last(n int) []Entry {
	h.mu.Lock()
	defer h.mu.Unlock()
	if n <= 0 || n >= len(h.entries) {
		copy := make([]Entry, len(h.entries))
		for i, e := range h.entries {
			copy[i] = e
		}
		return copy
	}
	slice := h.entries[len(h.entries)-n:]
	copy := make([]Entry, len(slice))
	for i, e := range slice {
		copy[i] = e
	}
	return copy
}
