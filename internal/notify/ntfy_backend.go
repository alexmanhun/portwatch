package notify

import (
	"fmt"
	"net/http"
	"strings"
)

// NtfyBackend sends notifications via ntfy.sh or a self-hosted ntfy server.
type NtfyBackend struct {
	topicURL string
	client   *http.Client
}

// NewNtfyBackend creates a NtfyBackend. topicURL is e.g. https://ntfy.sh/mytopic.
func NewNtfyBackend(topicURL string) *NtfyBackend {
	return &NtfyBackend{topicURL: topicURL, client: &http.Client{}}
}

func (n *NtfyBackend) Name() string { return "ntfy" }

func (n *NtfyBackend) Send(event Event) error {
	req, err := http.NewRequest(http.MethodPost, n.topicURL, strings.NewReader(event.Message))
	if err != nil {
		return err
	}
	req.Header.Set("Title", fmt.Sprintf("portwatch: %s", event.Type))
	req.Header.Set("Content-Type", "text/plain")

	resp, err := n.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("ntfy: unexpected status %d", resp.StatusCode)
	}
	return nil
}
