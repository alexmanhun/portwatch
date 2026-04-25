package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/alert"
)

// GooglePubSubBackend publishes port events to a Google Cloud Pub/Sub topic
// via the REST API. It requires a valid OAuth2/service-account bearer token.
type GooglePubSubBackend struct {
	projectID string
	topicID   string
	token     string
	client    *http.Client
}

// NewGooglePubSubBackend creates a backend that publishes to the given
// project/topic using the supplied bearer token.
func NewGooglePubSubBackend(projectID, topicID, token string) *GooglePubSubBackend {
	return &GooglePubSubBackend{
		projectID: projectID,
		topicID:   topicID,
		token:     token,
		client:    &http.Client{},
	}
}

func (g *GooglePubSubBackend) Name() string { return "googlepubsub" }

func (g *GooglePubSubBackend) Send(event alert.Event) error {
	data, err := json.Marshal(map[string]interface{}{
		"port":      event.Port,
		"eventType": event.Type,
		"timestamp": event.Timestamp,
	})
	if err != nil {
		return fmt.Errorf("googlepubsub: marshal error: %w", err)
	}

	// Pub/Sub expects base64-encoded data; encoding/json handles []byte as base64.
	body, err := json.Marshal(map[string]interface{}{
		"messages": []map[string]interface{}{
			{"data": data},
		},
	})
	if err != nil {
		return fmt.Errorf("googlepubsub: envelope error: %w", err)
	}

	url := fmt.Sprintf(
		"https://pubsub.googleapis.com/v1/projects/%s/topics/%s:publish",
		g.projectID, g.topicID,
	)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("googlepubsub: request error: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+g.token)

	resp, err := g.client.Do(req)
	if err != nil {
		return fmt.Errorf("googlepubsub: http error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("googlepubsub: unexpected status %d", resp.StatusCode)
	}
	return nil
}
