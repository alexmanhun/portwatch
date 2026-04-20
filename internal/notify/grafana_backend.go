package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// GrafanaBackend sends annotations to a Grafana instance when port events occur.
type GrafanaBackend struct {
	url    string
	apiKey string
	client *http.Client
}

// NewGrafanaBackend creates a new GrafanaBackend.
// url should be the base URL of the Grafana instance (e.g. http://localhost:3000).
// apiKey should be a Grafana API key with editor permissions.
func NewGrafanaBackend(url, apiKey string) *GrafanaBackend {
	return &GrafanaBackend{
		url:    url,
		apiKey: apiKey,
		client: &http.Client{},
	}
}

func (g *GrafanaBackend) Name() string { return "grafana" }

func (g *GrafanaBackend) Send(event Event) error {
	type annotationPayload struct {
		Text string   `json:"text"`
		Tags []string `json:"tags"`
	}

	tags := []string{"portwatch", event.Type}
	payload := annotationPayload{
		Text: fmt.Sprintf("[portwatch] %s: port %d", event.Type, event.Port),
		Tags: tags,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("grafana: marshal payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, g.url+"/api/annotations", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("grafana: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+g.apiKey)

	resp, err := g.client.Do(req)
	if err != nil {
		return fmt.Errorf("grafana: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("grafana: unexpected status %d", resp.StatusCode)
	}
	return nil
}
