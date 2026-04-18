package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// GotifyBackend sends notifications to a Gotify server.
type GotifyBackend struct {
	url   string
	token string
	client *http.Client
}

// NewGotifyBackend creates a new GotifyBackend.
func NewGotifyBackend(serverURL, token string) *GotifyBackend {
	return &GotifyBackend{
		url:    serverURL,
		token:  token,
		client: &http.Client{},
	}
}

func (g *GotifyBackend) Name() string { return "gotify" }

func (g *GotifyBackend) Send(event Event) error {
	payload := map[string]interface{}{
		"title":    fmt.Sprintf("portwatch: %s", event.Type),
		"message":  event.Message,
		"priority": 5,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, g.url+"/message", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Gotify-Key", g.token)

	resp, err := g.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("gotify: unexpected status %d", resp.StatusCode)
	}
	return nil
}
