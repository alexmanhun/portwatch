package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/alert"
)

// GoogleChatBackend sends notifications to a Google Chat webhook.
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

func (b *GoogleChatBackend) Name() string { return "googlechat" }

func (b *GoogleChatBackend) Send(event alert.Event) error {
	payload := map[string]string{
		"text": fmt.Sprintf("[portwatch] %s: port %d", event.Type, event.Port),
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	resp, err := b.client.Post(b.webhookURL, "application/json", bytes.NewReader(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("googlechat: unexpected status %d", resp.StatusCode)
	}
	return nil
}
