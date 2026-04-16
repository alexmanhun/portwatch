package history

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

// Entry records a port state change event.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Port      int       `json:"port"`
	Event     string    `json:"event"` // "opened" or "closed"
}

// History manages a persistent log of port change events.
type History struct {
	mu      sync.Mutex
	entries []Entry
	path    string
}

// New creates a History backed by the given file path.
// Existing entries are loaded if the file exists.
func New(path string) (*History, error) {
	h := &History{path: path}
	if err := h.load(); err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	return h, nil
}

// Record appends a new entry and persists to disk.
func (h *History) Record(port int, event string) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.entries = append(h.entries, Entry{
		Timestamp: time.Now().UTC(),
		Port:      port,
		Event:     event,
	})
	return h.save()
}

// Entries returns a copy of all recorded entries.
func (h *History) Entries() []Entry {
	h.mu.Lock()
	defer h.mu.Unlock()
	copy := make([]Entry, len(h.entries))
	for i, e := range h.entries {
		copy[i] = e
	}
	return copy
}

func (h *History) load() error {
	data, err := os.ReadFile(h.path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &h.entries)
}

func (h *History) save() error {
	data, err := json.MarshalIndent(h.entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(h.path, data, 0644)
}
