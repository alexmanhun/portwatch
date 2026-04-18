package main

import (
	"github.com/user/portwatch/internal/notify"
	"github.com/user/portwatch/internal/config"
)

// buildDispatcher constructs a Dispatcher wired with all configured backends.
func buildDispatcher(cfg *config.Config) *notify.Dispatcher {
	d := notify.New()

	// Log backend is always active.
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

	if cfg.DiscordWebhookURL != "" {
		d.Add(notify.NewDiscordBackend(cfg.DiscordWebhookURL))
	}

	if cfg.SMSWebhookURL != "" {
		d.Add(notify.NewSMSBackend(cfg.SMSWebhookURL))
	}

	if cfg.EmailHost != "" && cfg.EmailFrom != "" && cfg.EmailTo != "" {
		d.Add(notify.NewEmailBackend(cfg.EmailHost, cfg.EmailFrom, cfg.EmailTo))
	}

	return d
}
