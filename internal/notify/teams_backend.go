package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// TeamsBackend sends notifications to a Microsoft Teams channel via webhook.
type TeamsBackend struct {
	webhookURL string
	client     *http.Client
}

// NewTeamsBackend creates a new TeamsBackend with the given webhook URL.
func NewTeamsBackend(webhookURL string) *TeamsBackend {
	return &TeamsBackend{
		webhookURL: webhookURL,
		client:     &http.Client{},
	}
}

func (t *TeamsBackend) Name() string {
	return "teams"
}

func (t *TeamsBackend) Send(event Event) error {
	payload := map[string]string{
		"@type":      "MessageCard",
		"@context":   "http://schema.org/extensions",
		"summary":    event.Message,
		"title":      fmt.Sprintf("Portwatch: %s", event.Kind),
		"text":       event.Message,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("teams: marshal error: %w", err)
	}

	resp, err := t.client.Post(t.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("teams: request error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("teams: unexpected status %d", resp.StatusCode)
	}
	return nil
}
