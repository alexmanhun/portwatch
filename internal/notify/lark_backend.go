package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/alert"
)

// LarkBackend sends notifications to a Lark (Feishu) incoming webhook.
type LarkBackend struct {
	webhookURL string
	client     *http.Client
}

// NewLarkBackend creates a new LarkBackend with the given webhook URL.
func NewLarkBackend(webhookURL string) *LarkBackend {
	return &LarkBackend{
		webhookURL: webhookURL,
		client:     &http.Client{},
	}
}

// Name returns the backend identifier.
func (b *LarkBackend) Name() string {
	return "lark"
}

// Send dispatches an alert event to the Lark webhook.
func (b *LarkBackend) Send(event alert.Event) error {
	payload := map[string]interface{}{
		"msg_type": "text",
		"content": map[string]string{
			"text": fmt.Sprintf("[portwatch] %s: port %d — %s", event.Type, event.Port, event.Message),
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("lark: marshal payload: %w", err)
	}

	resp, err := b.client.Post(b.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("lark: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("lark: unexpected status %d", resp.StatusCode)
	}

	return nil
}
