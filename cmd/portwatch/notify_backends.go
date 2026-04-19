package main

import (
	"portwatch/internal/config"
	"portwatch/internal/notify"
)

func buildDispatcher(cfg *config.Config) *notify.Dispatcher {
	d := notify.New()
	d.Add(notify.NewLogBackend())

	if cfg.WebhookURL != "" {
		d.Add(notify.NewWebhookBackend(cfg.WebhookURL))
	}
	if cfg.PagerDutyKey != "" {
		d.Add(notify.NewPagerDutyBackend(cfg.PagerDutyKey))
	}
	if cfg.SlackWebhookURL != "" {
		d.Add(notify.NewSlackBackend(cfg.SlackWebhookURL))
	}
	if cfg.DatadogAPIKey != "" {
		d.Add(notify.NewDatadogBackend(cfg.DatadogAPIKey))
	}
	if cfg.NewRelicAccountID != "" && cfg.NewRelicInsertKey != "" {
		d.Add(notify.NewNewRelicBackend(cfg.NewRelicAccountID, cfg.NewRelicInsertKey))
	}

	return d
}
