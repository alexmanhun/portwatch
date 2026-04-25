package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/alert"
)

// CustomEventBackend sends port events to a generic HTTP endpoint with
// a user-defined JSON template rendered from the event fields.
type CustomEventBackend struct {
	url     string
	headers map[string]string
	client  *http.Client
}

// NewCustomEventBackend creates a backend that POSTs to url with optional
// extra headers (e.g. Authorization tokens).
func NewCustomEventBackend(url string, headers map[string]string) *CustomEventBackend {
	if headers == nil {
		headers = map[string]string{}
	}
	return &CustomEventBackend{
		url:     url,
		headers: headers,
		client:  &http.Client{},
	}
}

func (b *CustomEventBackend) Name() string { return "customevent" }

func (b *CustomEventBackend) Send(event alert.Event) error {
	payload := map[string]interface{}{
		"event": event.Type,
		"port":  event.Port,
		"proto": event.Proto,
		"time":  event.Time.UTC().Format("2006-01-02T15:04:05Z"),
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("customevent: marshal: %w", err)
	}
	req, err := http.NewRequest(http.MethodPost, b.url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("customevent: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	for k, v := range b.headers {
		req.Header.Set(k, v)
	}
	resp, err := b.client.Do(req)
	if err != nil {
		return fmt.Errorf("customevent: send: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("customevent: unexpected status %d", resp.StatusCode)
	}
	return nil
}
