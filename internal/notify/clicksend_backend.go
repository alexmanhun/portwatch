package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// ClickSendBackend sends SMS alerts via the ClickSend REST API.
type ClickSendBackend struct {
	username string
	apiKey   string
	to       string
	client   *http.Client
}

// NewClickSendBackend creates a new ClickSend notification backend.
// username is the ClickSend account username, apiKey is the API key,
// and to is the destination phone number in E.164 format.
func NewClickSendBackend(username, apiKey, to string) *ClickSendBackend {
	return &ClickSendBackend{
		username: username,
		apiKey:   apiKey,
		to:       to,
		client:   &http.Client{},
	}
}

func (b *ClickSendBackend) Name() string { return "clicksend" }

func (b *ClickSendBackend) Send(event Event) error {
	const apiURL = "https://rest.clicksend.com/v3/sms/send"

	payload := map[string]interface{}{
		"messages": []map[string]string{
			{
				"source": "portwatch",
				"body":   fmt.Sprintf("[portwatch] %s: port %d", event.Type, event.Port),
				"to":     b.to,
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("clicksend: marshal payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, apiURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("clicksend: create request: %w", err)
	}
	req.SetBasicAuth(b.username, b.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := b.client.Do(req)
	if err != nil {
		return fmt.Errorf("clicksend: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("clicksend: unexpected status %d", resp.StatusCode)
	}
	return nil
}
