package main

import (
	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/notify"
)

// buildDispatcher constructs a Dispatcher wired with all configured backends.
func buildDispatcher(cfg *config.Config) *notify.Dispatcher {
	d := notify.New()
	d.Register(notify.NewLogBackend())

	if cfg.WebhookURL != "" {
		d.Register(notify.NewWebhookBackend(cfg.WebhookURL))
	}
	if cfg.PagerDutyKey != "" {
		d.Register(notify.NewPagerDutyBackend(cfg.PagerDutyKey))
	}
	if cfg.SlackWebhook != "" {
		d.Register(notify.NewSlackBackend(cfg.SlackWebhook))
	}
	if cfg.SplunkURL != "" && cfg.SplunkToken != "" {
		d.Register(notify.NewSplunkBackend(cfg.SplunkURL, cfg.SplunkToken))
	}

	return d
}
