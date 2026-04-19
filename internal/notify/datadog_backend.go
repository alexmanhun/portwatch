package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const datadogEventsURL = "https://api.datadoghq.com/api/v1/events"

type DatadogBackend struct {
	apiKey string
	url    string
}

func NewDatadogBackend(apiKey string) *DatadogBackend {
	return &DatadogBackend{apiKey: apiKey, url: datadogEventsURL}
}

func (d *DatadogBackend) Name() string { return "datadog" }

func (d *DatadogBackend) Send(event Event) error {
	payload := map[string]interface{}{
		"title": fmt.Sprintf("portwatch: %s", event.Type),
		"text":  fmt.Sprintf("Port %d %s", event.Port, event.Type),
		"tags":  []string{"portwatch", fmt.Sprintf("port:%d", event.Port)},
		"alert_type": func() string {
			if event.Type == "closed" {
				return "warning"
			}
			return "info"
		}(),
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, d.url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("DD-API-KEY", d.apiKey)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("datadog: unexpected status %d", resp.StatusCode)
	}
	return nil
}
