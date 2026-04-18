package main

import (
	"github.com/user/portwatch/internal/notify"
	"github.com/user/portwatch/internal/config"
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

	if cfg.DiscordWebhookURL != "" {
		d.Add(notify.NewDiscordBackend(cfg.DiscordWebhookURL))
	}

	if cfg.TeamsWebhookURL != "" {
		d.Add(notify.NewTeamsBackend(cfg.TeamsWebhookURL))
	}

	if cfg.SMSWebhookURL != "" {
		d.Add(notify.NewSMSBackend(cfg.SMSWebhookURL))
	}

	if cfg.OpsGenieAPIKey != "" {
		d.Add(notify.NewOpsGenieBackend(cfg.OpsGenieAPIKey))
	}

	if cfg.EmailHost != "" && cfg.EmailFrom != "" && cfg.EmailTo != "" {
		d.Add(notify.NewEmailBackend(cfg.EmailHost, cfg.EmailFrom, cfg.EmailTo))
	}

	return d
}
