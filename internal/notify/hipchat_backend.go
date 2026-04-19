package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// HipChatBackend sends alerts to a HipChat room via the v2 API.
type HipChatBackend struct {
	token  string
	roomID string
	baseURL string
}

// NewHipChatBackend creates a HipChatBackend.
// baseURL defaults to https://api.hipchat.com if empty.
func NewHipChatBackend(token, roomID, baseURL string) *HipChatBackend {
	if baseURL == "" {
		baseURL = "https://api.hipchat.com"
	}
	return &HipChatBackend{token: token, roomID: roomID, baseURL: baseURL}
}

func (h *HipChatBackend) Name() string { return "hipchat" }

func (h *HipChatBackend) Send(event Event) error {
	payload := map[string]interface{}{
		"message":        fmt.Sprintf("[portwatch] %s port %d", event.Type, event.Port),
		"message_format": "text",
		"color":          colorForEvent(event.Type),
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%s/v2/room/%s/notification?auth_token=%s", h.baseURL, h.roomID, h.token)
	resp, err := http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("hipchat: unexpected status %d", resp.StatusCode)
	}
	return nil
}

func colorForEvent(eventType string) string {
	switch eventType {
	case "new":
		return "yellow"
	case "closed":
		return "red"
	default:
		return "gray"
	}
}
