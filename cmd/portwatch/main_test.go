package main

import (
	"os"
	"testing"

	"portwatch/internal/config"
)

func TestConfigDefaultUsedWhenNoFile(t *testing.T) {
	cfg := config.Default()
	if cfg == nil {
		t.Fatal("expected non-nil default config")
	}
	if cfg.PortRange.Start < 1 {
		t.Errorf("invalid default start port: %d", cfg.PortRange.Start)
	}
	if cfg.PortRange.End < cfg.PortRange.Start {
		t.Errorf("default end port %d less than start port %d", cfg.PortRange.End, cfg.PortRange.Start)
	}
}

func TestAlertOutputFileCreation(t *testing.T) {
	tmp, err := os.CreateTemp("", "portwatch-alert-*.log")
	if err != nil {
		t.Fatal(err)
	}
	tmp.Close()
	defer os.Remove(tmp.Name())

	f, err := os.OpenFile(tmp.Name(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		t.Fatalf("failed to open alert output file: %v", err)
	}
	f.Close()

	if _, err := os.Stat(tmp.Name()); os.IsNotExist(err) {
		t.Error("alert output file was not created")
	}
}
