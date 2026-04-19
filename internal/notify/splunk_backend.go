package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/alert"
)

// SplunkBackend sends events to a Splunk HTTP Event Collector (HEC) endpoint.
type SplunkBackend struct {
	url   string
	token string
	client *http.Client
}

// NewSplunkBackend creates a new SplunkBackend.
// url should be the full HEC endpoint, e.g. https://splunk:8088/services/collector/event
func NewSplunkBackend(url, token string) *SplunkBackend {
	return &SplunkBackend{url: url, token: token, client: &http.Client{}}
}

func (s *SplunkBackend) Name() string { return "splunk" }

func (s *SplunkBackend) Send(e alert.Event) error {
	payload := map[string]interface{}{
		"event": map[string]interface{}{
			"port":    e.Port,
			"type":    e.Type,
			"message": e.Message,
		},
		"sourcetype": "portwatch",
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("splunk: marshal: %w", err)
	}
	req, err := http.NewRequest(http.MethodPost, s.url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("splunk: request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Splunk "+s.token)
	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("splunk: send: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("splunk: unexpected status %d", resp.StatusCode)
	}
	return nil
}
