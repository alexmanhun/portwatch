package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// KafkaBackend sends events to a Kafka REST Proxy endpoint.
type KafkaBackend struct {
	proxyURL string
	topic    string
	client   *http.Client
}

// NewKafkaBackend creates a backend targeting the given Kafka REST Proxy URL and topic.
func NewKafkaBackend(proxyURL, topic string) *KafkaBackend {
	return &KafkaBackend{
		proxyURL: proxyURL,
		topic:    topic,
		client:   &http.Client{Timeout: 5 * time.Second},
	}
}

func (k *KafkaBackend) Name() string { return "kafka" }

func (k *KafkaBackend) Send(event Event) error {
	record, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("kafka: marshal: %w", err)
	}

	payload, _ := json.Marshal(map[string]any{
		"records": []map[string]any{
			{"value": json.RawMessage(record)},
		},
	})

	url := fmt.Sprintf("%s/topics/%s", k.proxyURL, k.topic)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("kafka: request: %w", err)
	}
	req.Header.Set("Content-Type", "application/vnd.kafka.json.v2+json")

	resp, err := k.client.Do(req)
	if err != nil {
		return fmt.Errorf("kafka: send: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("kafka: unexpected status %d", resp.StatusCode)
	}
	return nil
}
