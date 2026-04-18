package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// SlackBackend sends alert notifications to a Slack webhook URL.
type SlackBackend struct {
	webhookURL string
	client     *http.Client
}

// NewSlackBackend creates a new SlackBackend with the given webhook URL.
func NewSlackBackend(webhookURL string) *SlackBackend {
	return &SlackBackend{
		webhookURL: webhookURL,
		client:     &http.Client{Timeout: 10 * time.Second},
	}
}

func (s *SlackBackend) Name() string {
	return "slack"
}

func (s *SlackBackend) Send(event Event) error {
	text := fmt.Sprintf("*portwatch alert* [%s] port %d — %s",
		event.Type, event.Port, event.Message)

	payload, err := json.Marshal(map[string]string{"text": text})
	if err != nil {
		return fmt.Errorf("slack: marshal error: %w", err)
	}

	resp, err := s.client.Post(s.webhookURL, "application/json", bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("slack: request error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("slack: unexpected status %d", resp.StatusCode)
	}
	return nil
}
