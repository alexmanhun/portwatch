package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/alert"
)

const pushbulletURL = "https://api.pushbullet.com/v2/pushes"

// PushbulletBackend sends notifications via the Pushbullet API.
type PushbulletBackend struct {
	apiKey string
	url    string
	client *http.Client
}

// NewPushbulletBackend creates a new PushbulletBackend with the given API key.
func NewPushbulletBackend(apiKey string) *PushbulletBackend {
	return &PushbulletBackend{
		apiKey: apiKey,
		url:    pushbulletURL,
		client: &http.Client{},
	}
}

func (p *PushbulletBackend) Name() string { return "pushbullet" }

func (p *PushbulletBackend) Send(event alert.Event) error {
	payload := map[string]interface{}{
		"type":  "note",
		"title": "portwatch alert",
		"body":  fmt.Sprintf("%s: port %d", event.Type, event.Port),
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("pushbullet: marshal: %w", err)
	}
	req, err := http.NewRequest(http.MethodPost, p.url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("pushbullet: new request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Access-Token", p.apiKey)

	resp, err := p.client.Do(req)
	if err != nil {
		return fmt.Errorf("pushbullet: do: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("pushbullet: unexpected status %d", resp.StatusCode)
	}
	return nil
}
