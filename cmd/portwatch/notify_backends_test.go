package main

import (
	"portwatch/internal/config"
	"testing"
)

func backendNames(d interface{ Names() []string }) []string {
	return d.Names()
}

func hasBackend(names []string, target string) bool {
	for _, n := range names {
		if n == target {
			return true
		}
	}
	return false
}

func TestBuildDispatcherAlwaysHasLogBackend(t *testing.T) {
	cfg := config.Default()
	d := buildDispatcher(cfg)
	if !hasBackend(d.Names(), "log") {
		t.Error("expected log backend to always be present")
	}
}

func TestBuildDispatcherWebhook(t *testing.T) {
	cfg := config.Default()
	cfg.WebhookURL = "http://example.com/hook"
	d := buildDispatcher(cfg)
	if !hasBackend(d.Names(), "webhook") {
		t.Error("expected webhook backend")
	}
}

func TestBuildDispatcherPagerDuty(t *testing.T) {
	cfg := config.Default()
	cfg.PagerDutyKey = "somekey"
	d := buildDispatcher(cfg)
	if !hasBackend(d.Names(), "pagerduty") {
		t.Error("expected pagerduty backend")
	}
}

func TestBuildDispatcherOpsGenie(t *testing.T) {
	cfg := config.Default()
	cfg.OpsGenieKey = "ogkey"
	d := buildDispatcher(cfg)
	if !hasBackend(d.Names(), "opsgenie") {
		t.Error("expected opsgenie backend")
	}
}

func TestBuildDispatcherDingTalk(t *testing.T) {
	cfg := config.Default()
	cfg.DingTalkWebhook = "http://oapi.dingtalk.com/robot/send?access_token=abc"
	d := buildDispatcher(cfg)
	if !hasBackend(d.Names(), "dingtalk") {
		t.Error("expected dingtalk backend")
	}
}
