package main

import (
	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/notify"
)

// buildDispatcher constructs a Dispatcher wired with all configured backends.
func buildDispatcher(cfg *config.Config) *notify.Dispatcher {
	d := notify.New()

	// Log backend is always present.
	d.Register(notify.NewLogBackend(nil))

	if cfg.WebhookURL != "" {
		d.Register(notify.NewWebhookBackend(cfg.WebhookURL))
	}
	if cfg.PagerDutyKey != "" {
		d.Register(notify.NewPagerDutyBackend(cfg.PagerDutyKey))
	}
	if cfg.SlackWebhook != "" {
		d.Register(notify.NewSlackBackend(cfg.SlackWebhook))
	}
	if cfg.EmailAddr != "" {
		d.Register(notify.NewEmailBackend(cfg.EmailAddr, cfg.SMTPHost))
	}
	if cfg.WebexToken != "" && cfg.WebexRoomID != "" {
		d.Register(notify.NewWebexBackend(cfg.WebexToken, cfg.WebexRoomID))
	}

	return d
}
