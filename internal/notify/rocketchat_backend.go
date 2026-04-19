package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// RocketChatBackend sends alerts to a Rocket.Chat incoming webhook.
type RocketChatBackend struct {
	webhookURL string
	client     *http.Client
}

// NewRocketChatBackend creates a new RocketChatBackend.
func NewRocketChatBackend(webhookURL string) *RocketChatBackend {
	return &RocketChatBackend{
		webhookURL: webhookURL,
		client:     &http.Client{},
	}
}

func (r *RocketChatBackend) Name() string {
	return "rocketchat"
}

func (r *RocketChatBackend) Send(event Event) error {
	payload := map[string]string{
		"text": fmt.Sprintf("[portwatch] %s: port %d", event.Type, event.Port),
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	resp, err := r.client.Post(r.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("rocketchat: unexpected status %d", resp.StatusCode)
	}
	return nil
}
