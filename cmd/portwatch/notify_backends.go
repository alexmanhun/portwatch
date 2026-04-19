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
	if cfg.PagerDutyKey != "" {
		d.Add(notify.NewPagerDutyBackend(cfg.PagerDutyKey))
	}
	if cfg.OpsGenieKey != "" {
		d.Add(notify.NewOpsGenieBackend(cfg.OpsGenieKey))
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
	if cfg.VictorOpsURL != "" {
		d.Add(notify.NewVictorOpsBackend(cfg.VictorOpsURL))
	}
	if cfg.GotifyURL != "" && cfg.GotifyToken != "" {
		d.Add(notify.NewGotifyBackend(cfg.GotifyURL, cfg.GotifyToken))
	}
	if cfg.NtfyURL != "" {
		d.Add(notify.NewNtfyBackend(cfg.NtfyURL))
	}
	if cfg.XMPPGatewayURL != "" && cfg.XMPPTo != "" {
		d.Add(notify.NewXMPPBackend(cfg.XMPPGatewayURL, cfg.XMPPTo))
	}

	return d
}
