package main

import (
	"testing"

	"github.com/user/portwatch/internal/config"
)

func backendNames(d interface{ Names() []string }) []string {
	return d.Names()
}

func hasBackend(names []string, name string) bool {
	for _, n := range names {
		if n == name {
			return true
		}
	}
	return false
}

func TestBuildDispatcherAlwaysHasLogBackend(t *testing.T) {
	cfg := config.Default()
	d := buildDispatcher(cfg)
	if !hasBackend(d.Names(), "log") {
		t.Error("expected log backend to always be registered")
	}
}

func TestBuildDispatcherWebhook(t *testing.T) {
	cfg := config.Default()
	cfg.WebhookURL = "http://example.com/hook"
	d := buildDispatcher(cfg)
	if !hasBackend(d.Names(), "webhook") {
		t.Error("expected webhook backend when WebhookURL is set")
	}
}

func TestBuildDispatcherPagerDuty(t *testing.T) {
	cfg := config.Default()
	cfg.PagerDutyKey = "pd-key-123"
	d := buildDispatcher(cfg)
	if !hasBackend(d.Names(), "pagerduty") {
		t.Error("expected pagerduty backend when PagerDutyKey is set")
	}
}

func TestBuildDispatcherWebex(t *testing.T) {
	cfg := config.Default()
	cfg.WebexToken = "wx-token"
	cfg.WebexRoomID = "room-id"
	d := buildDispatcher(cfg)
	if !hasBackend(d.Names(), "webex") {
		t.Error("expected webex backend when WebexToken and WebexRoomID are set")
	}
}

func TestBuildDispatcherWebexSkippedWhenPartial(t *testing.T) {
	cfg := config.Default()
	cfg.WebexToken = "wx-token"
	// WebexRoomID intentionally left empty
	d := buildDispatcher(cfg)
	if hasBackend(d.Names(), "webex") {
		t.Error("expected webex backend to be skipped when RoomID is missing")
	}
}
