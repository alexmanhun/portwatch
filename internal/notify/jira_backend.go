package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/alert"
)

// JiraBackend creates a Jira issue when a port event occurs.
type JiraBackend struct {
	baseURL   string
	project   string
	apiToken  string
	username  string
	client    *http.Client
}

// NewJiraBackend constructs a JiraBackend.
// baseURL is the Jira instance root (e.g. https://myorg.atlassian.net).
func NewJiraBackend(baseURL, project, username, apiToken string) *JiraBackend {
	return &JiraBackend{
		baseURL:  baseURL,
		project:  project,
		apiToken: apiToken,
		username: username,
		client:   &http.Client{},
	}
}

func (j *JiraBackend) Name() string { return "jira" }

func (j *JiraBackend) Send(event alert.Event) error {
	summary := fmt.Sprintf("[portwatch] %s on port %d", event.Type, event.Port)
	description := fmt.Sprintf("Port %d changed state: %s at %s", event.Port, event.Type, event.Timestamp.Format("2006-01-02T15:04:05Z"))

	body := map[string]interface{}{
		"fields": map[string]interface{}{
			"project":     map[string]string{"key": j.project},
			"summary":     summary,
			"description": description,
			"issuetype":   map[string]string{"name": "Bug"},
		},
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("jira: marshal: %w", err)
	}

	url := j.baseURL + "/rest/api/2/issue"
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("jira: request: %w", err)
	}
	req.SetBasicAuth(j.username, j.apiToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := j.client.Do(req)
	if err != nil {
		return fmt.Errorf("jira: send: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("jira: unexpected status %d", resp.StatusCode)
	}
	return nil
}
