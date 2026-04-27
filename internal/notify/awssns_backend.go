package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/alert"
)

type awsSNSBackend struct {
	topicARN  string
	region    string
	accessKey string
	secretKey string
	endpoint  string
	client    *http.Client
}

// NewAWSSNSBackend creates a Backend that publishes alerts to an AWS SNS topic.
// For production use, credentials should be provided via environment variables
// or IAM roles; the accessKey/secretKey parameters are used for explicit auth.
func NewAWSSNSBackend(topicARN, region, accessKey, secretKey string) Backend {
	return &awsSNSBackend{
		topicARN:  topicARN,
		region:    region,
		accessKey: accessKey,
		secretKey: secretKey,
		endpoint:  fmt.Sprintf("https://sns.%s.amazonaws.com", region),
		client:    &http.Client{},
	}
}

func (b *awsSNSBackend) Name() string { return "awssns" }

func (b *awsSNSBackend) Send(ev alert.Event) error {
	message := fmt.Sprintf("portwatch alert: %s on port %d", ev.Type, ev.Port)

	payload := map[string]string{
		"TopicArn": b.topicARN,
		"Message":  message,
		"Subject":  fmt.Sprintf("portwatch [%s] port %d", ev.Type, ev.Port),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("awssns: marshal payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, b.endpoint, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("awssns: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if b.accessKey != "" {
		req.SetBasicAuth(b.accessKey, b.secretKey)
	}

	resp, err := b.client.Do(req)
	if err != nil {
		return fmt.Errorf("awssns: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("awssns: unexpected status %d", resp.StatusCode)
	}
	return nil
}
