package config

import (
	"encoding/json"
	"os"
)

// Config holds all portwatch configuration.
type Config struct {
	PortRange     string `json:"port_range"`
	Interval      int    `json:"interval_seconds"`
	AlertOutput   string `json:"alert_output"`
	HistoryFile   string `json:"history_file"`
	WebhookURL    string `json:"webhook_url,omitempty"`
	SlackURL      string `json:"slack_url,omitempty"`
	DiscordURL    string `json:"discord_url,omitempty"`
	TeamsURL      string `json:"teams_url,omitempty"`
	PagerDutyKey  string `json:"pagerduty_key,omitempty"`
	OpsGenieKey   string `json:"opsgenie_key,omitempty"`
	GoogleChatURL string `json:"googlechat_url,omitempty"`
}

// Default returns a Config populated with sensible defaults.
func Default() *Config {
	return &Config{
		PortRange:   "1-1024",
		Interval:    60,
		AlertOutput: "",
		HistoryFile: "portwatch_history.json",
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

// Save writes a Config to a JSON file.
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
