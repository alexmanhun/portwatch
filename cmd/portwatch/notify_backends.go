package main

import (
	"portwatch/internal/notify"
	"portwatch/internal/config"
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
	if cfg.OpsGenieKey != "" {
		d.Register(notify.NewOpsGenieBackend(cfg.OpsGenieKey))
	}
	if cfg.SlackWebhook != "" {
		d.Register(notify.NewSlackBackend(cfg.SlackWebhook))
	}
	if cfg.DingTalkWebhook != "" {
		d.Register(notify.NewDingTalkBackend(cfg.DingTalkWebhook))
	}

	return d
}
