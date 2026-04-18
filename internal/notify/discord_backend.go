package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// DiscordBackend sends alerts to a Discord webhook.
type DiscordBackend struct {
	webhookURL string
	client     *http.Client
}

// NewDiscordBackend creates a new DiscordBackend with the given webhook URL.
func NewDiscordBackend(webhookURL string) *DiscordBackend {
	return &DiscordBackend{
		webhookURL: webhookURL,
		client:     &http.Client{},
	}
}

func (d *DiscordBackend) Name() string { return "discord" }

func (d *DiscordBackend) Send(event Event) error {
	payload := map[string]string{
		"content": fmt.Sprintf("**portwatch** [%s] port %d — %s", event.Type, event.Port, event.Message),
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("discord: marshal: %w", err)
	}
	resp, err := d.client.Post(d.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("discord: post: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("discord: unexpected status %d", resp.StatusCode)
	}
	return nil
}
