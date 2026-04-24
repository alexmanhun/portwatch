package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/alert"
)

// BearyChatBackend sends notifications to a BearyChat incoming webhook.
type BearyChatBackend struct {
	webhookURL string
	client     *http.Client
}

// NewBearyChatBackend creates a new BearyChatBackend with the given webhook URL.
func NewBearyChatBackend(webhookURL string) *BearyChatBackend {
	return &BearyChatBackend{
		webhookURL: webhookURL,
		client:     &http.Client{},
	}
}

// Name returns the backend identifier.
func (b *BearyChatBackend) Name() string {
	return "bearychat"
}

// bearyChatPayload represents the JSON body sent to BearyChat.
type bearyChatPayload struct {
	Text        string `json:"text"`
	Notification string `json:"notification"`
}

// Send dispatches the event to the BearyChat webhook.
func (b *BearyChatBackend) Send(event alert.Event) error {
	payload := bearyChatPayload{
		Text:        fmt.Sprintf("**portwatch**: port %d — %s", event.Port, event.Type),
		Notification: fmt.Sprintf("portwatch: port %d %s", event.Port, event.Type),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("bearychat: marshal: %w", err)
	}

	resp, err := b.client.Post(b.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("bearychat: post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("bearychat: unexpected status %d", resp.StatusCode)
	}
	return nil
}
