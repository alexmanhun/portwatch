package main

import (
	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/notify"
)

// buildDispatcher constructs a Dispatcher wired up with all configured
// notification backends. The log backend is always included.
func buildDispatcher(cfg *config.Config) *notify.Dispatcher {
	d := notify.New()
	d.Register(notify.NewLogBackend())

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
		d.Register(notify.NewEmailBackend(cfg.EmailSMTP, cfg.EmailAddr))
	}
	if cfg.TwilioAccountSID != "" && cfg.TwilioAuthToken != "" {
		d.Register(notify.NewTwilioBackend(
			cfg.TwilioAccountSID,
			cfg.TwilioAuthToken,
			cfg.TwilioFrom,
			cfg.TwilioTo,
		))
	}

	return d
}
