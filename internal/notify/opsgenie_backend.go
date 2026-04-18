package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const opsGenieAPIURL = "https://api.opsgenie.com/v2/alerts"

type OpsGenieBackend struct {
	apiKey  string
	apiURL  string
	client  *http.Client
}

func NewOpsGenieBackend(apiKey string) *OpsGenieBackend {
	return &OpsGenieBackend{
		apiKey: apiKey,
		apiURL: opsGenieAPIURL,
		client: &http.Client{},
	}
}

func (o *OpsGenieBackend) Name() string {
	return "opsgenie"
}

func (o *OpsGenieBackend) Send(event Event) error {
	payload := map[string]interface{}{
		"message":     fmt.Sprintf("portwatch: %s on port %d", event.Type, event.Port),
		"description": event.Detail,
		"priority":    "P3",
		"tags":        []string{"portwatch", string(event.Type)},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, o.apiURL, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "GenieKey "+o.apiKey)

	resp, err := o.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("opsgenie: unexpected status %d", resp.StatusCode)
	}
	return nil
}
