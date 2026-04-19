package main

import (
	"github.com/user/portwatch/internal/notify"
	"github.com/user/portwatch/internal/config"
)

// buildDispatcher constructs a Dispatcher wired with all configured backends.
func buildDispatcher(cfg *config.Config) *notify.Dispatcher {
	d := notify.New()
	d.Add(notify.NewLogBackend())

	if cfg.WebhookURL != "" {
		d.Add(notify.NewWebhookBackend(cfg.WebhookURL))
	}
	if cfg.SlackURL != "" {
		d.Add(notify.NewSlackBackend(cfg.SlackURL))
	}
	if cfg.DiscordURL != "" {
		d.Add(notify.NewDiscordBackend(cfg.DiscordURL))
	}
	if cfg.TeamsURL != "" {
		d.Add(notify.NewTeamsBackend(cfg.TeamsURL))
	}
	if cfg.PagerDutyKey != "" {
		d.Add(notify.NewPagerDutyBackend(cfg.PagerDutyKey))
	}
	if cfg.OpsGenieKey != "" {
		d.Add(notify.NewOpsGenieBackend(cfg.OpsGenieKey))
	}
	if cfg.GoogleChatURL != "" {
		d.Add(notify.NewGoogleChatBackend(cfg.GoogleChatURL))
	}
	return d
}
