package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/alert"
)

// DingTalkBackend sends notifications via DingTalk webhook.
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

func (b *DingTalkBackend) Name() string { return "dingtalk" }

func (b *DingTalkBackend) Send(event alert.Event) error {
	payload := map[string]interface{}{
		"msgtype": "text",
		"text": map[string]string{
			"content": fmt.Sprintf("[portwatch] %s: port %d", event.Type, event.Port),
		},
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	resp, err := b.client.Post(b.webhookURL, "application/json", bytes.NewReader(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("dingtalk: unexpected status %d", resp.StatusCode)
	}
	return nil
}
