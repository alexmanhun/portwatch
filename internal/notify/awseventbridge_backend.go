package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// AWSEventBridgeBackend sends port events to an AWS EventBridge partner event
// bus via the PutEvents HTTP endpoint (or a compatible local proxy).
type AWSEventBridgeBackend struct {
	endpointURL string
	source      string
	detailType  string
	apiKey      string
	client      *http.Client
}

type eventBridgeEntry struct {
	Source     string          `json:"Source"`
	DetailType string          `json:"DetailType"`
	Detail     json.RawMessage `json:"Detail"`
	Time       string          `json:"Time"`
}

type eventBridgeRequest struct {
	Entries []eventBridgeEntry `json:"Entries"`
}

// NewAWSEventBridgeBackend creates a backend that forwards events to the given
// EventBridge HTTP endpoint. source and detailType follow the EventBridge
// schema; apiKey is sent as the X-Api-Key header (optional).
func NewAWSEventBridgeBackend(endpointURL, source, detailType, apiKey string) *AWSEventBridgeBackend {
	if source == "" {
		source = "portwatch"
	}
	if detailType == "" {
		detailType = "PortEvent"
	}
	return &AWSEventBridgeBackend{
		endpointURL: endpointURL,
		source:      source,
		detailType:  detailType,
		apiKey:      apiKey,
		client:      &http.Client{Timeout: 10 * time.Second},
	}
}

func (b *AWSEventBridgeBackend) Name() string { return "awseventbridge" }

func (b *AWSEventBridgeBackend) Send(event Event) error {
	detail, err := json.Marshal(map[string]interface{}{
		"port":      event.Port,
		"eventType": event.Type,
		"message":   event.Message,
	})
	if err != nil {
		return fmt.Errorf("awseventbridge: marshal detail: %w", err)
	}

	body, err := json.Marshal(eventBridgeRequest{
		Entries: []eventBridgeEntry{
			{
				Source:     b.source,
				DetailType: b.detailType,
				Detail:     detail,
				Time:       time.Now().UTC().Format(time.RFC3339),
			},
		},
	})
	if err != nil {
		return fmt.Errorf("awseventbridge: marshal request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, b.endpointURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("awseventbridge: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if b.apiKey != "" {
		req.Header.Set("X-Api-Key", b.apiKey)
	}

	resp, err := b.client.Do(req)
	if err != nil {
		return fmt.Errorf("awseventbridge: send: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("awseventbridge: unexpected status %d", resp.StatusCode)
	}
	return nil
}
