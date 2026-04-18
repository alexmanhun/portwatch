package main

import (
	"github.com/user/portwatch/internal/notify"
	"github.com/user/portwatch/internal/config"
)

// buildDispatcher constructs a Dispatcher wired with all configured backends.
func buildDispatcher(cfg *config.Config) *notify.Dispatcher {
	d := notify.New()

	// Log backend is always active
	d.Register(notify.NewLogBackend())

	if cfg.WebhookURL != "" {
		d.Register(notify.NewWebhookBackend(cfg.WebhookURL))
	}
	if cfg.SlackWebhookURL != "" {
		d.Register(notify.NewSlackBackend(cfg.SlackWebhookURL))
	}
	if cfg.DiscordWebhookURL != "" {
		d.Register(notify.NewDiscordBackend(cfg.DiscordWebhookURL))
	}
	if cfg.TeamsWebhookURL != "" {
		d.Register(notify.NewTeamsBackend(cfg.TeamsWebhookURL))
	}
	if cfg.PagerDutyKey != "" {
		d.Register(notify.NewPagerDutyBackend(cfg.PagerDutyKey))
	}
	if cfg.OpsGenieKey != "" {
		d.Register(notify.NewOpsGenieBackend(cfg.OpsGenieKey))
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
	if cfg.MatrixURL != "" && cfg.MatrixToken != "" {
		d.Register(notify.NewMatrixBackend(cfg.MatrixURL, cfg.MatrixToken))
	}
	if cfg.TelegramBotToken != "" && cfg.TelegramChatID != "" {
		d.Register(notify.NewTelegramBackend(cfg.TelegramBotToken, cfg.TelegramChatID))
	}

	return d
}
