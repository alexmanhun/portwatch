package main

import (
	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/notify"
)

// buildDispatcher constructs a Dispatcher wired with every backend that has
// been configured. The log backend is always present as a fallback.
func buildDispatcher(cfg *config.Config) *notify.Dispatcher {
	d := notify.New()
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
	if cfg.CustomEventURL != "" {
		headers := map[string]string{}
		if cfg.CustomEventToken != "" {
			headers["Authorization"] = "Bearer " + cfg.CustomEventToken
		}
		d.Register(notify.NewCustomEventBackend(cfg.CustomEventURL, headers))
	}

	return d
}
