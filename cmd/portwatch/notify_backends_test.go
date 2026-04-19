package main

import (
	"portwatch/internal/config"
	"testing"
)

func backendNames(d interface{ Backends() []string }) []string {
	return d.Backends()
}

func TestBuildDispatcherAlwaysHasLogBackend(t *testing.T) {
	cfg := config.Default()
	d := buildDispatcher(cfg)
	for _, name := range d.Backends() {
		if name == "log" {
			return
		}
	}
	t.Fatal("log backend missing")
}

func TestBuildDispatcherWebhook(t *testing.T) {
	cfg := config.Default()
	cfg.WebhookURL = "http://example.com/hook"
	d := buildDispatcher(cfg)
	for _, name := range d.Backends() {
		if name == "webhook" {
			return
		}
	}
	t.Fatal("webhook backend missing")
}

func TestBuildDispatcherPagerDuty(t *testing.T) {
	cfg := config.Default()
	cfg.PagerDutyKey = "key123"
	d := buildDispatcher(cfg)
	for _, name := range d.Backends() {
		if name == "pagerduty" {
			return
		}
	}
	t.Fatal("pagerduty backend missing")
}

func TestBuildDispatcherOpsGenie(t *testing.T) {
	cfg := config.Default()
	cfg.OpsGenieKey = "ogkey"
	d := buildDispatcher(cfg)
	for _, name := range d.Backends() {
		if name == "opsgenie" {
			return
		}
	}
	t.Fatal("opsgenie backend missing")
}

func TestBuildDispatcherMQTT(t *testing.T) {
	cfg := config.Default()
	cfg.MQTTBrokerURL = "http://localhost:8080/publish"
	d := buildDispatcher(cfg)
	for _, name := range d.Backends() {
		if name == "mqtt" {
			return
		}
	}
	t.Fatal("mqtt backend missing")
}
