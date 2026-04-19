package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// SignalRBackend sends notifications to an Azure SignalR REST endpoint.
type SignalRBackend struct {
	endpoint string
	hubName  string
	apiKey   string
	client   *http.Client
}

// NewSignalRBackend creates a new SignalRBackend.
func NewSignalRBackend(endpoint, hubName, apiKey string) *SignalRBackend {
	return &SignalRBackend{
		endpoint: endpoint,
		hubName:  hubName,
		apiKey:   apiKey,
		client:   &http.Client{},
	}
}

func (s *SignalRBackend) Name() string { return "signalr" }

func (s *SignalRBackend) Send(event Event) error {
	url := fmt.Sprintf("%s/api/v1/hubs/%s", s.endpoint, s.hubName)

	payload := map[string]interface{}{
		"target":    "portwatch",
		"arguments": []interface{}{event.Type, event.Port, event.Message},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("signalr: marshal error: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("signalr: request error: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if s.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+s.apiKey)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("signalr: send error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("signalr: unexpected status %d", resp.StatusCode)
	}
	return nil
}
