package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// WebhookBackend posts JSON notifications to an HTTP endpoint.
type WebhookBackend struct {
	url    string
	client *http.Client
}

// NewWebhookBackend creates a WebhookBackend targeting url.
func NewWebhookBackend(url string) *WebhookBackend {
	return &WebhookBackend{
		url:    url,
		client: &http.Client{Timeout: 5 * time.Second},
	}
}

func (w *WebhookBackend) Name() string { return "webhook" }

type webhookPayload struct {
	Level string `json:"level"`
	Title string `json:"title"`
	Body  string `json:"body"`
	Port  int    `json:"port"`
	Event string `json:"event"`
}

func (w *WebhookBackend) Send(msg Message) error {
	payload := webhookPayload{
		Level: string(msg.Level),
		Title: msg.Title,
		Body:  msg.Body,
		Port:  msg.Port,
		Event: msg.Event,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	resp, err := w.client.Post(w.url, "application/json", bytes.NewReader(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("webhook returned status %d", resp.StatusCode)
	}
	return nil
}
