package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const newRelicEventsURL = "https://insights-collector.newrelic.com/v1/accounts/%s/events"

type NewRelicBackend struct {
	accountID string
	insertKey string
	url       string
}

func NewNewRelicBackend(accountID, insertKey string) *NewRelicBackend {
	return &NewRelicBackend{
		accountID: accountID,
		insertKey: insertKey,
		url:       fmt.Sprintf(newRelicEventsURL, accountID),
	}
}

func (n *NewRelicBackend) Name() string { return "newrelic" }

func (n *NewRelicBackend) Send(event Event) error {
	payload := []map[string]interface{}{
		{
			"eventType": "PortwatchEvent",
			"port":      event.Port,
			"eventKind": event.Type,
		},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, n.url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Insert-Key", n.insertKey)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("newrelic: unexpected status %d", resp.StatusCode)
	}
	return nil
}
