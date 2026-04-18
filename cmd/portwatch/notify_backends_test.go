package main

import (
	"testing"

	"portwatch/internal/config"
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
	t.Fatal("expected log backend to always be present")
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
	t.Fatal("expected webhook backend")
}

func TestBuildDispatcherPagerDuty(t *testing.T) {
	cfg := config.Default()
	cfg.PagerDutyKey = "somekey"
	d := buildDispatcher(cfg)
	for _, name := range d.Backends() {
		if name == "pagerduty" {
			return
		}
	}
	t.Fatal("expected pagerduty backend")
}

func TestBuildDispatcherNoEmailWithoutAllFields(t *testing.T) {
	cfg := config.Default()
	cfg.EmailHost = "smtp.example.com"
	// missing From and To
	d := buildDispatcher(cfg)
	for _, name := range d.Backends() {
		if name == "email" {
			t.Fatal("email backend should not be added without all fields")
		}
	}
}

func TestBuildDispatcherSMS(t *testing.T) {
	cfg := config.Default()
	cfg.SMSGatewayURL = "http://sms.example.com"
	cfg.SMSAPIKey = "key"
	cfg.SMSFrom = "+1000"
	cfg.SMSTo = "+2000"
	d := buildDispatcher(cfg)
	for _, name := range d.Backends() {
		if name == "sms" {
			return
		}
	}
	t.Fatal("expected sms backend")
}
