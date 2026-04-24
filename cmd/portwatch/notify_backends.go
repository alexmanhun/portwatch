package main

import (
	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/notify"
)

// buildDispatcher constructs a Dispatcher wired with all backends that are
// enabled in cfg. The log backend is always included.
func buildDispatcher(cfg *config.Config) *notify.Dispatcher {
	d := notify.New()

	d.Register(notify.NewLogBackend())

	if cfg.WebhookURL != "" {
		d.Register(notify.NewWebhookBackend(cfg.WebhookURL))
	}

	if cfg.PagerDutyKey != "" {
		d.Register(notify.NewPagerDutyBackend(cfg.PagerDutyKey))
	}

	if cfg.SlackWebhookURL != "" {
		d.Register(notify.NewSlackBackend(cfg.SlackWebhookURL))
	}

	if cfg.EmailAddr != "" {
		d.Register(notify.NewEmailBackend(cfg.EmailAddr, cfg.SMTPHost))
	}

	if cfg.FlowdockToken != "" {
		d.Register(notify.NewFlowdockBackend(cfg.FlowdockToken))
	}

	return d
}
