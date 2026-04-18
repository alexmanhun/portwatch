package config

import (
	"encoding/json"
	"os"
)

// Config holds all portwatch runtime configuration.
type Config struct {
	ScanRange    string `json:"scan_range"`
	IntervalSecs int    `json:"interval_secs"`
	AlertOutput  string `json:"alert_output"`
	HistoryFile  string `json:"history_file"`

	// Webhook
	WebhookURL string `json:"webhook_url"`

	// PagerDuty
	PagerDutyKey string `json:"pagerduty_key"`

	// Slack
	SlackWebhookURL string `json:"slack_webhook_url"`

	// Email
	EmailHost string `json:"email_host"`
	EmailFrom string `json:"email_from"`
	EmailTo   string `json:"email_to"`

	// SMS
	SMSGatewayURL string `json:"sms_gateway_url"`
	SMSAPIKey     string `json:"sms_api_key"`
	SMSFrom       string `json:"sms_from"`
	SMSTo         string `json:"sms_to"`
}

// Default returns a Config populated with sensible defaults.
func Default() *Config {
	return &Config{
		ScanRange:    "1-1024",
		IntervalSecs: 60,
		AlertOutput:  "",
		HistoryFile:  "portwatch_history.json",
	}
}

// Load reads a Config from a JSON file.
func Load(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	cfg := Default()
	if err := json.NewDecoder(f).Decode(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

// Save writes cfg to a JSON file at path.
func Save(cfg *Config, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(cfg)
}
