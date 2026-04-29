package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/alert"
)

// ClickSendBackend sends SMS notifications via the ClickSend API.
type ClickSendBackend struct {
	apiURL   string
	username string
	apiKey   string
	to       string
	client   *http.Client
}

// NewClickSendBackend creates a new ClickSendBackend.
func NewClickSendBackend(apiURL, username, apiKey, to string) *ClickSendBackend {
	return &ClickSendBackend{
		apiURL:   apiURL,
		username: username,
		apiKey:   apiKey,
		to:       to,
		client:   &http.Client{},
	}
}

func (c *ClickSendBackend) Name() string { return "clicksend" }

func (c *ClickSendBackend) Send(event alert.Event) error {
	payload := map[string]interface{}{
		"messages": []map[string]string{
			{
				"to":   c.to,
				"body": fmt.Sprintf("[portwatch] %s port %d: %s", event.Type, event.Port, event.Message),
			},
		},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("clicksend: marshal: %w", err)
	}
	req, err := http.NewRequest(http.MethodPost, c.apiURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("clicksend: request: %w", err)
	}
	req.SetBasicAuth(c.username, c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("clicksend: do: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("clicksend: unexpected status %d", resp.StatusCode)
	}
	return nil
}
