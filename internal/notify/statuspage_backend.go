package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/alert"
)

// StatusPageBackend sends incident updates to a Atlassian Statuspage API.
type StatusPageBackend struct {
	apiKey    string
	pageID    string
	componentID string
	baseURL   string
	client    *http.Client
}

// NewStatusPageBackend creates a backend that posts component status updates
// to Atlassian Statuspage. baseURL defaults to the public API if empty.
func NewStatusPageBackend(apiKey, pageID, componentID, baseURL string) *StatusPageBackend {
	if baseURL == "" {
		baseURL = "https://api.statuspage.io/v1"
	}
	return &StatusPageBackend{
		apiKey:      apiKey,
		pageID:      pageID,
		componentID: componentID,
		baseURL:     baseURL,
		client:      &http.Client{},
	}
}

func (s *StatusPageBackend) Name() string { return "statuspage" }

func (s *StatusPageBackend) Send(e alert.Event) error {
	status := "operational"
	if e.Type == alert.EventNewPort {
		status = "under_maintenance"
	}

	payload := map[string]interface{}{
		"component": map[string]string{
			"status": status,
		},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("statuspage: marshal: %w", err)
	}

	url := fmt.Sprintf("%s/pages/%s/components/%s", s.baseURL, s.pageID, s.componentID)
	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("statuspage: build request: %w", err)
	}
	req.Header.Set("Authorization", "OAuth "+s.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("statuspage: request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("statuspage: unexpected status %d", resp.StatusCode)
	}
	return nil
}
