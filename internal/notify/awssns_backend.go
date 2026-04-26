package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// AWSSNSBackend publishes port-change events to an AWS SNS topic via the
// SNS HTTP publish endpoint. It is intended for use with SNS HTTP/HTTPS
// subscriptions or a local SNS-compatible mock (e.g. LocalStack).
type AWSSNSBackend struct {
	topicARN string
	endpointURL string // e.g. https://sns.us-east-1.amazonaws.com or LocalStack URL
	accessKey string
	secretKey string
	region string
	client *http.Client
}

// NewAWSSNSBackend creates a new AWSSNSBackend.
// endpointURL is the base SNS endpoint (leave empty to use the default AWS
// regional endpoint derived from region).
func NewAWSSNSBackend(topicARN, region, accessKey, secretKey, endpointURL string) *AWSSNSBackend {
	ep := endpointURL
	if ep == "" {
		ep = fmt.Sprintf("https://sns.%s.amazonaws.com", region)
	}
	return &AWSSNSBackend{
		topicARN:    topicARN,
		endpointURL: ep,
		accessKey:   accessKey,
		secretKey:   secretKey,
		region:      region,
		client:      &http.Client{Timeout: 10 * time.Second},
	}
}

// Name returns the backend identifier.
func (b *AWSSNSBackend) Name() string { return "awssns" }

// snsPublishPayload is the JSON body sent to the SNS publish endpoint.
// This mirrors the structure expected by LocalStack and SNS-compatible APIs
// that accept JSON rather than query-string form encoding.
type snsPublishPayload struct {
	TopicArn string `json:"TopicArn"`
	Message  string `json:"Message"`
	Subject  string `json:"Subject"`
}

// Send dispatches the event to the configured SNS topic.
func (b *AWSSNSBackend) Send(event Event) error {
	subject := fmt.Sprintf("portwatch: %s on port %d", event.Type, event.Port)
	message := fmt.Sprintf(
		"{\"event\":%q,\"port\":%d,\"timestamp\":%q}",
		event.Type, event.Port, event.Timestamp.UTC().Format(time.RFC3339),
	)

	payload := snsPublishPayload{
		TopicArn: b.topicARN,
		Message:  message,
		Subject:  subject,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("awssns: marshal payload: %w", err)
	}

	// Build the publish URL. Real AWS uses query-params + SigV4; for
	// simplicity (and LocalStack compatibility) we POST JSON to the endpoint.
	url := b.endpointURL + "/?Action=Publish"

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("awssns: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if b.accessKey != "" {
		// Minimal auth header — production deployments should use proper SigV4.
		req.Header.Set("X-Amz-Access-Key", b.accessKey)
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
