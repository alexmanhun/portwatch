package main

import (
	"testing"

	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/notify"
)

func backendNames(d *notify.Dispatcher) []string {
	names := make([]string, 0)
	for _, b := range d.Backends() {
		names = append(names, b.Name())
	}
	return names
}

func hasBackend(d *notify.Dispatcher, name string) bool {
	for _, n := range backendNames(d) {
		if n == name {
			return true
		}
	}
	return false
}

func TestBuildDispatcherAlwaysHasLogBackend(t *testing.T) {
	cfg := config.Default()
	d := buildDispatcher(cfg)
	if !hasBackend(d, "log") {
		t.Error("expected log backend to always be present")
	}
}

func TestBuildDispatcherWebhook(t *testing.T) {
	cfg := config.Default()
	cfg.WebhookURL = "https://example.com/hook"
	d := buildDispatcher(cfg)
	if !hasBackend(d, "webhook") {
		t.Error("expected webhook backend when WebhookURL is set")
	}
}

func TestBuildDispatcherPagerDuty(t *testing.T) {
	cfg := config.Default()
	cfg.PagerDutyKey = "somekey"
	d := buildDispatcher(cfg)
	if !hasBackend(d, "pagerduty") {
		t.Error("expected pagerduty backend when PagerDutyKey is set")
	}
}

func TestBuildDispatcherTwilio(t *testing.T) {
	cfg := config.Default()
	cfg.TwilioAccountSID = "ACtest"
	cfg.TwilioAuthToken = "token"
	cfg.TwilioFrom = "+10000000000"
	cfg.TwilioTo = "+19999999999"
	d := buildDispatcher(cfg)
	if !hasBackend(d, "twilio") {
		t.Error("expected twilio backend when Twilio credentials are set")
	}
}

func TestBuildDispatcherTwilioMissingToken(t *testing.T) {
	cfg := config.Default()
	cfg.TwilioAccountSID = "ACtest"
	// AuthToken intentionally omitted
	d := buildDispatcher(cfg)
	if hasBackend(d, "twilio") {
		t.Error("twilio backend should not be registered without auth token")
	}
}
