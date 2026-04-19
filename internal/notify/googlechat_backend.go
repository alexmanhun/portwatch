package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// GoogleChatBackend sends alerts to a Google Chat webhook.
type GoogleChatBackend struct {
	webhookURL string
	client     *http.Client
}

// NewGoogleChatBackend creates a new GoogleChatBackend.
func NewGoogleChatBackend(webhookURL string) *GoogleChatBackend {
	return &GoogleChatBackend{
		webhookURL: webhookURL,
		client:     &http.Client{},
	}
}

func (g *GoogleChatBackend) Name() string { return "googlechat" }

func (g *GoogleChatBackend) Send(event Event) error {
	payload := map[string]string{
		"text": fmt.Sprintf("*portwatch* [%s] port %d: %s", event.Type, event.Port, event.Message),
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	resp, err := g.client.Post(g.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("googlechat: unexpected status %d", resp.StatusCode)
	}
	return nil
}
