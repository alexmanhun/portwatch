package main

import (
	"os"
	"testing"
)

func TestBuildDispatcherAlwaysHasLogBackend(t *testing.T) {
	d := buildDispatcher()
	names := backendNames(d)
	if !contains(names, "log") {
		t.Error("expected log backend to always be present")
	}
}

func TestBuildDispatcherWebhook(t *testing.T) {
	os.Setenv("PORTWATCH_WEBHOOK_URL", "http://example.com/hook")
	defer os.Unsetenv("PORTWATCH_WEBHOOK_URL")

	d := buildDispatcher()
	if !contains(backendNames(d), "webhook") {
		t.Error("expected webhook backend when env var is set")
	}
}

func TestBuildDispatcherPagerDuty(t *testing.T) {
	os.Setenv("PORTWATCH_PAGERDUTY_KEY", "somekey")
	defer os.Unsetenv("PORTWATCH_PAGERDUTY_KEY")

	d := buildDispatcher()
	if !contains(backendNames(d), "pagerduty") {
		t.Error("expected pagerduty backend when env var is set")
	}
}

func TestBuildDispatcherNoEmailWithoutAllFields(t *testing.T) {
	os.Setenv("PORTWATCH_SMTP_HOST", "smtp.example.com")
	defer os.Unsetenv("PORTWATCH_SMTP_HOST")
	// Missing PORTWATCH_ALERT_EMAIL and PORTWATCH_FROM_EMAIL

	d := buildDispatcher()
	if contains(backendNames(d), "email") {
		t.Error("email backend should not be registered without to/from addresses")
	}
}

// backendNames returns the names of all registered backends via a test dispatch.
func backendNames(d interface{ RegisteredNames() []string }) []string {
	return d.RegisteredNames()
}

func contains(ss []string, s string) bool {
	for _, v := range ss {
		if v == s {
			return true
		}
	}
	return false
}
