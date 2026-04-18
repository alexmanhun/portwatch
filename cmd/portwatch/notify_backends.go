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
	if cfg.OpsGenieKey != "" {
		d.Register(notify.NewOpsGenieBackend(cfg.OpsGenieKey))
	}
	if cfg.SlackWebhook != "" {
		d.Register(notify.NewSlackBackend(cfg.SlackWebhook))
	}
	if cfg.DiscordWebhook != "" {
		d.Register(notify.NewDiscordBackend(cfg.DiscordWebhook))
	}
	if cfg.TeamsWebhook != "" {
		d.Register(notify.NewTeamsBackend(cfg.TeamsWebhook))
	}
	if cfg.VictorOpsURL != "" {
		d.Register(notify.NewVictorOpsBackend(cfg.VictorOpsURL))
	}
	if cfg.GotifyURL != "" && cfg.GotifyToken != "" {
		d.Register(notify.NewGotifyBackend(cfg.GotifyURL, cfg.GotifyToken))
	}
	if cfg.NtfyURL != "" {
		d.Register(notify.NewNtfyBackend(cfg.NtfyURL))
	}
	if cfg.PushoverToken != "" && cfg.PushoverUser != "" {
		d.Register(notify.NewPushoverBackend(cfg.PushoverToken, cfg.PushoverUser))
	}
	return d
}
