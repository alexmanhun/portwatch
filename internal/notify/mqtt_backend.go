package notify

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// MQTTBackend publishes alerts to an MQTT broker via HTTP bridge (e.g. EMQX REST API).
type MQTTBackend struct {
	brokerURL string
	topic     string
	client    *http.Client
}

// NewMQTTBackend creates a new MQTTBackend.
// brokerURL should be the HTTP API endpoint, e.g. http://localhost:8080/api/v5/publish
func NewMQTTBackend(brokerURL, topic string) *MQTTBackend {
	return &MQTTBackend{
		brokerURL: brokerURL,
		topic:     topic,
		client:    &http.Client{Timeout: 5 * time.Second},
	}
}

func (m *MQTTBackend) Name() string { return "mqtt" }

func (m *MQTTBackend) Send(event Event) error {
	payload := map[string]interface{}{
		"topic":   m.topic,
		"qos":     1,
		"payload": fmt.Sprintf("%s port %d", event.Type, event.Port),
	}
	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	resp, err := m.client.Post(m.brokerURL, "application/json", strings.NewReader(string(b)))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("mqtt: unexpected status %d", resp.StatusCode)
	}
	return nil
}
