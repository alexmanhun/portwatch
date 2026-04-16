package config

import (
	"encoding/json"
	"os"
	"time"
)

// Config holds the portwatch daemon configuration.
type Config struct {
	PortRange   PortRange     `json:"port_range"`
	Interval    time.Duration `json:"interval"`
	AlertOutput string        `json:"alert_output"`
	Baseline    []int         `json:"baseline"`
}

// PortRange defines the start and end ports to scan.
type PortRange struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

// Default returns a Config populated with sensible defaults.
func Default() *Config {
	return &Config{
		PortRange:   PortRange{Start: 1, End: 1024},
		Interval:    30 * time.Second,
		AlertOutput: "",
		Baseline:    []int{},
	}
}

// Load reads a JSON config file from the given path.
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

// Save writes the config as JSON to the given path.
func Save(path string, cfg *Config) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(cfg)
}
