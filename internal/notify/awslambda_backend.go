package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// AWSLambdaBackend invokes an AWS Lambda function URL with port event details.
type AWSLambdaBackend struct {
	functionURL string
	apiKey      string
	client      *http.Client
}

// NewAWSLambdaBackend creates a backend that calls an AWS Lambda function URL.
// functionURL is the HTTPS endpoint exposed by the Lambda function URL feature.
// apiKey is sent as the x-api-key header; pass an empty string to omit it.
func NewAWSLambdaBackend(functionURL, apiKey string) *AWSLambdaBackend {
	return &AWSLambdaBackend{
		functionURL: functionURL,
		apiKey:      apiKey,
		client:      &http.Client{},
	}
}

func (b *AWSLambdaBackend) Name() string { return "awslambda" }

func (b *AWSLambdaBackend) Send(event Event) error {
	payload := map[string]interface{}{
		"source":  "portwatch",
		"event":   event.Type,
		"port":    event.Port,
		"message": event.Message,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("awslambda: marshal: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, b.functionURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("awslambda: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if b.apiKey != "" {
		req.Header.Set("x-api-key", b.apiKey)
	}

	resp, err := b.client.Do(req)
	if err != nil {
		return fmt.Errorf("awslambda: send: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("awslambda: unexpected status %d", resp.StatusCode)
	}
	return nil
}
