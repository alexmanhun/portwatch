package notify

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const pushoverAPIURL = "https://api.pushover.net/1/messages.json"

// PushoverBackend sends notifications via the Pushover API.
type PushoverBackend struct {
	token   string
	userKey string
	apiURL  string
	client  *http.Client
}

// NewPushoverBackend creates a new PushoverBackend.
func NewPushoverBackend(token, userKey string) *PushoverBackend {
	return &PushoverBackend{
		token:   token,
		userKey: userKey,
		apiURL:  pushoverAPIURL,
		client:  &http.Client{},
	}
}

func (p *PushoverBackend) Name() string { return "pushover" }

func (p *PushoverBackend) Send(event Event) error {
	payload := map[string]string{
		"token":   p.token,
		"user":    p.userKey,
		"title":   fmt.Sprintf("portwatch: %s", event.Type),
		"message": event.Message,
	}
	b, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("pushover: marshal: %w", err)
	}
	resp, err := p.client.Post(p.apiURL, "application/json", strings.NewReader(string(b)))
	if err != nil {
		return fmt.Errorf("pushover: request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("pushover: unexpected status %d", resp.StatusCode)
	}
	return nil
}
