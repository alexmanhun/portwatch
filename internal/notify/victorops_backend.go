package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type VictorOpsBackend struct {
	webhookURL string
	client     *http.Client
}

func NewVictorOpsBackend(webhookURL string) *VictorOpsBackend {
	return &VictorOpsBackend{
		webhookURL: webhookURL,
		client:     &http.Client{},
	}
}

func (v *VictorOpsBackend) Name() string {
	return "victorops"
}

func (v *VictorOpsBackend) Send(event Event) error {
	msgType := "INFO"
	if event.Kind == "new_port" {
		msgType = "WARNING"
	} else if event.Kind == "closed_port" {
		msgType = "INFO"
	}

	payload := map[string]any{
		"message_type":  msgType,
		"entity_id":     fmt.Sprintf("portwatch-port-%d", event.Port),
		"state_message": event.Message,
		"port":          event.Port,
		"kind":          event.Kind,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("victorops: marshal error: %w", err)
	}

	resp, err := v.client.Post(v.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("victorops: request error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("victorops: unexpected status %d", resp.StatusCode)
	}
	return nil
}
