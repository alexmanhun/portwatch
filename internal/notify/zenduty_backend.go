package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// ZendutyBackend sends alerts to Zenduty via their Incidents API.
type ZendutyBackend struct {
	apiKey      string
	serviceID   string
	escalationID string
	client      *http.Client
}

// NewZendutyBackend creates a new ZendutyBackend.
func NewZendutyBackend(apiKey, serviceID, escalationID string) *ZendutyBackend {
	return &ZendutyBackend{
		apiKey:       apiKey,
		serviceID:    serviceID,
		escalationID: escalationID,
		client:       &http.Client{},
	}
}

func (z *ZendutyBackend) Name() string { return "zenduty" }

func (z *ZendutyBackend) Send(event Event) error {
	payload := map[string]interface{}{
		"title":         fmt.Sprintf("portwatch: %s on port %d", event.Type, event.Port),
		"message":       event.Message,
		"service":       z.serviceID,
		"escalation_policy": z.escalationID,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("zenduty: marshal: %w", err)
	}
	req, err := http.NewRequest(http.MethodPost,
		"https://www.zenduty.com/api/incidents/", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("zenduty: request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Token "+z.apiKey)

	resp, err := z.client.Do(req)
	if err != nil {
		return fmt.Errorf("zenduty: send: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("zenduty: unexpected status %d", resp.StatusCode)
	}
	return nil
}
