package main

import (
	"testing"

	"github.com/user/portwatch/internal/config"
)

func backendNames(d interface{ Backends() []string }) []string {
	return d.Backends()
}

func TestBuildDispatcherAlwaysHasLogBackend(t *testing.T) {
	cfg := config.Default()
	d := buildDispatcher(cfg)
	for _, n := range d.Backends() {
		if n == "log" {
			return
		}
	}
	t.Fatal("log backend missing")
}

func TestBuildDispatcherWebhook(t *testing.T) {
	cfg := config.Default()
	cfg.WebhookURL = "http://example.com/hook"
	d := buildDispatcher(cfg)
	for _, n := range d.Backends() {
		if n == "webhook" {
			return
		}
	}
	t.Fatal("webhook backend missing")
}

func TestBuildDispatcherPagerDuty(t *testing.T) {
	cfg := config.Default()
	cfg.PagerDutyKey = "testkey"
	d := buildDispatcher(cfg)
	for _, n := range d.Backends() {
		if n == "pagerduty" {
			return
		}
	}
	t.Fatal("pagerduty backend missing")
}

func TestBuildDispatcherOpsGenie(t *testing.T) {
	cfg := config.Default()
	cfg.OpsGenieKey = "ogkey"
	d := buildDispatcher(cfg)
	for _, n := range d.Backends() {
		if n == "opsgenie" {
			return
		}
	}
	t.Fatal("opsgenie backend missing")
}

func TestBuildDispatcherKafka(t *testing.T) {
	cfg := config.Default()
	cfg.KafkaProxyURL = "http://localhost:8082"
	cfg.KafkaTopic = "portwatch"
	d := buildDispatcher(cfg)
	for _, n := range d.Backends() {
		if n == "kafka" {
			return
		}
	}
	t.Fatal("kafka backend missing")
}
