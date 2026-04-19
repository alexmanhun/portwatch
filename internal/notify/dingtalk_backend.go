package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// DingTalkBackend sends alerts to a DingTalk webhook.
type DingTalkBackend struct {
	webhookURL string
	client     *http.Client
}

// NewDingTalkBackend creates a new DingTalkBackend.
func NewDingTalkBackend(webhookURL string) *DingTalkBackend {
	return &DingTalkBackend{
		webhookURL: webhookURL,
		client:     &http.Client{},
	}
}

func (d *DingTalkBackend) Name() string { return "dingtalk" }

func (d *DingTalkBackend) Send(event Event) error {
	payload := map[string]interface{}{
		"msgtype": "text",
		"text": map[string]string{
			"content": fmt.Sprintf("[portwatch] %s port %d", event.Type, event.Port),
		},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	resp, err := d.client.Post(d.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("dingtalk: unexpected status %d", resp.StatusCode)
	}
	return nil
}
