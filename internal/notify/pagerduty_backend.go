package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const pagerDutyEventURL = "https://events.pagerduty.com/v2/enqueue"

// PagerDutyBackend sends alerts to PagerDuty via Events API v2.
type PagerDutyBackend struct {
	integrationKey string
	client         *http.Client
	url            string
}

// NewPagerDutyBackend creates a new PagerDutyBackend with the given integration key.
func NewPagerDutyBackend(integrationKey string) *PagerDutyBackend {
	return &PagerDutyBackend{
		integrationKey: integrationKey,
		client:         &http.Client{},
		url:            pagerDutyEventURL,
	}
}

func (p *PagerDutyBackend) Name() string { return "pagerduty" }

func (p *PagerDutyBackend) Send(event Event) error {
	payload := map[string]interface{}{
		"routing_key":  p.integrationKey,
		"event_action": "trigger",
		"payload": map[string]interface{}{
			"summary":  event.Message,
			"severity": "warning",
			"source":   "portwatch",
			"custom_details": map[string]interface{}{
				"port": event.Port,
				"type": event.Type,
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("pagerduty: marshal error: %w", err)
	}

	resp, err := p.client.Post(p.url, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("pagerduty: request error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("pagerduty: unexpected status %d", resp.StatusCode)
	}
	return nil
}
