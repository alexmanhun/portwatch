package main

import (
	"portwatch/internal/config"
	"portwatch/internal/notify"
)

// buildDispatcher constructs a Dispatcher wired with all configured backends.
func buildDispatcher(cfg *config.Config) *notify.Dispatcher {
	d := notify.New()

	// Always include the log backend
	d.Add(notify.NewLogBackend())

	if cfg.WebhookURL != "" {
		d.Add(notify.NewWebhookBackend(cfg.WebhookURL))
	}
	if cfg.PagerDutyKey != "" {
		d.Add(notify.NewPagerDutyBackend(cfg.PagerDutyKey))
	}
	if cfg.OpsGenieKey != "" {
		d.Add(notify.NewOpsGenieBackend(cfg.OpsGenieKey))
	}
	if cfg.SlackWebhook != "" {
		d.Add(notify.NewSlackBackend(cfg.SlackWebhook))
	}
	if cfg.SignalREndpoint != "" && cfg.SignalRHub != "" {
		d.Add(notify.NewSignalRBackend(cfg.SignalREndpoint, cfg.SignalRHub, cfg.SignalRKey))
	}

	return d
}
