package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// MattermostBackend sends alerts to a Mattermost incoming webhook.
type MattermostBackend struct {
	webhookURL string
	client     *http.Client
}

// NewMattermostBackend creates a new MattermostBackend.
func NewMattermostBackend(webhookURL string) *MattermostBackend {
	return &MattermostBackend{
		webhookURL: webhookURL,
		client:     &http.Client{},
	}
}

func (m *MattermostBackend) Name() string {
	return "mattermost"
}

func (m *MattermostBackend) Send(event Event) error {
	payload := map[string]string{
		"text": fmt.Sprintf("**portwatch**: %s port %d at %s",
			event.Type, event.Port, event.Time.Format("2006-01-02 15:04:05")),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("mattermost: marshal error: %w", err)
	}

	resp, err := m.client.Post(m.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("mattermost: request error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("mattermost: unexpected status %d", resp.StatusCode)
	}
	return nil
}
