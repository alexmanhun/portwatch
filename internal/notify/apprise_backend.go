package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/alert"
)

// AppriseBackend sends notifications via an Apprise API server.
// Apprise supports 80+ notification services through a single unified API.
type AppriseBackend struct {
	serverURL string
	urls      []string
	client    *http.Client
}

// NewAppriseBackend creates a Dispatcher backend that posts to an Apprise
// API server at serverURL, forwarding the event to each of the provided
// notification URLs (e.g. "slack://token/channel", "mailto://user@host").
func NewAppriseBackend(serverURL string, urls []string) *AppriseBackend {
	return &AppriseBackend{
		serverURL: serverURL,
		urls:      urls,
		client:    &http.Client{},
	}
}

func (a *AppriseBackend) Name() string { return "apprise" }

func (a *AppriseBackend) Send(event alert.Event) error {
	payload := map[string]interface{}{
		"title": fmt.Sprintf("portwatch: %s", event.Type),
		"body":  fmt.Sprintf("Port %d — %s", event.Port, event.Type),
		"urls":  a.urls,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("apprise: marshal payload: %w", err)
	}

	endpoint := a.serverURL + "/notify"
	resp, err := a.client.Post(endpoint, "application/json", bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("apprise: post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("apprise: unexpected status %d", resp.StatusCode)
	}
	return nil
}
