package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// RedisBackend publishes alert events to a Redis channel via the HTTP REST API
// (e.g. Upstash Redis or a custom redis-http bridge).
type RedisBackend struct {
	url     string
	channel string
	token   string
	client  *http.Client
}

// NewRedisBackend creates a RedisBackend that publishes to the given channel.
func NewRedisBackend(url, channel, token string) *RedisBackend {
	return &RedisBackend{
		url:     url,
		channel: channel,
		token:   token,
		client:  &http.Client{Timeout: 5 * time.Second},
	}
}

func (r *RedisBackend) Name() string { return "redis" }

func (r *RedisBackend) Send(event string, port int) error {
	payload := map[string]interface{}{
		"channel": r.channel,
		"message": fmt.Sprintf("{\"event\":%q,\"port\":%d}", event, port),
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("redis: marshal: %w", err)
	}
	req, err := http.NewRequest(http.MethodPost, r.url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("redis: request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if r.token != "" {
		req.Header.Set("Authorization", "Bearer "+r.token)
	}
	resp, err := r.client.Do(req)
	if err != nil {
		return fmt.Errorf("redis: send: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("redis: unexpected status %d", resp.StatusCode)
	}
	return nil
}
