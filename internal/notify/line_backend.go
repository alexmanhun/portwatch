package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/alert"
)

// LineBackend sends notifications to a LINE Notify endpoint.
type LineBackend struct {
	token   string
	endpoint string
	client  *http.Client
}

// NewLineBackend creates a new LINE Notify backend.
// token is the LINE Notify personal access token.
func NewLineBackend(token string) *LineBackend {
	return &LineBackend{
		token:    token,
		endpoint: "https://notify-api.line.me/api/notify",
		client:   &http.Client{},
	}
}

func (b *LineBackend) Name() string { return "line" }

func (b *LineBackend) Send(event alert.Event) error {
	payload := map[string]string{
		"message": fmt.Sprintf("[portwatch] %s port %d: %s", event.Type, event.Port, event.Message),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("line: marshal payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, b.endpoint, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("line: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+b.token)

	resp, err := b.client.Do(req)
	if err != nil {
		return fmt.Errorf("line: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("line: unexpected status %d", resp.StatusCode)
	}
	return nil
}
