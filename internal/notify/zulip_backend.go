package notify

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// ZulipBackend sends alerts to a Zulip stream via the Zulip REST API.
type ZulipBackend struct {
	baseURL  string
	email    string
	apiKey   string
	stream   string
	topic    string
	client   *http.Client
}

// NewZulipBackend creates a new ZulipBackend.
// baseURL is the Zulip server URL (e.g. https://yourorg.zulipchat.com).
func NewZulipBackend(baseURL, email, apiKey, stream, topic string) *ZulipBackend {
	return &ZulipBackend{
		baseURL: strings.TrimRight(baseURL, "/"),
		email:   email,
		apiKey:  apiKey,
		stream:  stream,
		topic:   topic,
		client:  &http.Client{},
	}
}

// Name returns the backend identifier.
func (z *ZulipBackend) Name() string { return "zulip" }

// Send dispatches an event notification to the configured Zulip stream.
func (z *ZulipBackend) Send(event Event) error {
	endpoint := fmt.Sprintf("%s/api/v1/messages", z.baseURL)

	message := fmt.Sprintf("**portwatch** — %s: port %d", event.Type, event.Port)
	if event.Detail != "" {
		message += fmt.Sprintf(" (%s)", event.Detail)
	}

	form := url.Values{}
	form.Set("type", "stream")
	form.Set("to", z.stream)
	form.Set("topic", z.topic)
	form.Set("content", message)

	req, err := http.NewRequest(http.MethodPost, endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return fmt.Errorf("zulip: build request: %w", err)
	}
	req.SetBasicAuth(z.email, z.apiKey)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := z.client.Do(req)
	if err != nil {
		return fmt.Errorf("zulip: send: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("zulip: unexpected status %d", resp.StatusCode)
	}
	return nil
}
