package main

import (
	"portwatch/internal/notify"

	"portwatch/internal/config"
)

// buildDispatcher constructs a Dispatcher from the provided config,
// always including the log backend and adding optional backends when configured.
func buildDispatcher(cfg *config.Config) *notify.Dispatcher {
	d := notify.New()
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
	if cfg.IRCServer != "" && cfg.IRCNick != "" && cfg.IRCChannel != "" {
		d.Add(notify.NewIRCBackend(cfg.IRCServer, cfg.IRCNick, cfg.IRCChannel))
	}

	return d
}
