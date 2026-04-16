package config_test

import (
	"os"
	"testing"
	"time"

	"portwatch/internal/config"
)

func TestDefault(t *testing.T) {
	cfg := config.Default()
	if cfg.PortRange.Start != 1 || cfg.PortRange.End != 1024 {
		t.Errorf("unexpected default port range: %+v", cfg.PortRange)
	}
	if cfg.Interval != 30*time.Second {
		t.Errorf("unexpected default interval: %v", cfg.Interval)
	}
}

func TestSaveAndLoad(t *testing.T) {
	tmp, err := os.CreateTemp("", "portwatch-config-*.json")
	if err != nil {
		t.Fatal(err)
	}
	tmp.Close()
	defer os.Remove(tmp.Name())

	orig := &config.Config{
		PortRange:   config.PortRange{Start: 8000, End: 9000},
		Interval:    60 * time.Second,
		AlertOutput: "/tmp/alerts.log",
		Baseline:    []int{8080, 8443},
	}

	if err := config.Save(tmp.Name(), orig); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	loaded, err := config.Load(tmp.Name())
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if loaded.PortRange != orig.PortRange {
		t.Errorf("port range mismatch: got %+v want %+v", loaded.PortRange, orig.PortRange)
	}
	if loaded.Interval != orig.Interval {
		t.Errorf("interval mismatch: got %v want %v", loaded.Interval, orig.Interval)
	}
	if loaded.AlertOutput != orig.AlertOutput {
		t.Errorf("alert output mismatch: got %q want %q", loaded.AlertOutput, orig.AlertOutput)
	}
	if len(loaded.Baseline) != len(orig.Baseline) {
		t.Errorf("baseline length mismatch: got %d want %d", len(loaded.Baseline), len(orig.Baseline))
	}
}

func TestLoadMissingFile(t *testing.T) {
	_, err := config.Load("/nonexistent/path/config.json")
	if err == nil {
		t.Error("expected error loading missing file, got nil")
	}
}
