package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/alert"
)

// WebexBackend sends alerts to a Cisco Webex room via the Webex REST API.
type WebexBackend struct {
	token  string
	roomID string
	url    string
	client *http.Client
}

// NewWebexBackend creates a new WebexBackend.
// token is the Webex bot token and roomID is the destination room.
func NewWebexBackend(token, roomID string) *WebexBackend {
	return &WebexBackend{
		token:  token,
		roomID: roomID,
		url:    "https://webexapis.com/v1/messages",
		client: &http.Client{},
	}
}

func (w *WebexBackend) Name() string { return "webex" }

func (w *WebexBackend) Send(event alert.Event) error {
	payload := map[string]string{
		"roomId": w.roomID,
		"text":   fmt.Sprintf("[portwatch] %s port %d: %s", event.Type, event.Port, event.Message),
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("webex: marshal payload: %w", err)
	}
	req, err := http.NewRequest(http.MethodPost, w.url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("webex: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+w.token)

	resp, err := w.client.Do(req)
	if err != nil {
		return fmt.Errorf("webex: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webex: unexpected status %d", resp.StatusCode)
	}
	return nil
}
