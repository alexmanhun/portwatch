package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/alert"
)

// FlowdockBackend sends alert notifications to a Flowdock flow via the
// Flowdock API using a flow token.
type FlowdockBackend struct {
	flowToken string
	client    *http.Client
}

// NewFlowdockBackend creates a new FlowdockBackend with the given flow token.
func NewFlowdockBackend(flowToken string) *FlowdockBackend {
	return &FlowdockBackend{
		flowToken: flowToken,
		client:    &http.Client{},
	}
}

func (b *FlowdockBackend) Name() string { return "flowdock" }

func (b *FlowdockBackend) Send(event alert.Event) error {
	url := fmt.Sprintf("https://api.flowdock.com/messages/chat/%s", b.flowToken)

	payload := map[string]interface{}{
		"content":         fmt.Sprintf("[portwatch] %s", event.Message),
		"external_user_name": "portwatch",
		"tags":            []string{"portwatch", string(event.Type)},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("flowdock: marshal error: %w", err)
	}

	resp, err := b.client.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("flowdock: request error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("flowdock: unexpected status %d", resp.StatusCode)
	}

	return nil
}
