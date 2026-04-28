package main

import (
	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/notify"
)

// buildDispatcher constructs a Dispatcher wired with all configured backends.
func buildDispatcher(cfg *config.Config) *notify.Dispatcher {
	d := notify.New()

	// Log backend is always active.
	d.Register(notify.NewLogBackend())

	if cfg.WebhookURL != "" {
		d.Register(notify.NewWebhookBackend(cfg.WebhookURL))
	}
	if cfg.SlackWebhookURL != "" {
		d.Register(notify.NewSlackBackend(cfg.SlackWebhookURL))
	}
	if cfg.PagerDutyKey != "" {
		d.Register(notify.NewPagerDutyBackend(cfg.PagerDutyKey))
	}
	if cfg.EmailSMTP != "" {
		d.Register(notify.NewEmailBackend(cfg.EmailSMTP, cfg.EmailFrom, cfg.EmailTo))
	}
	if cfg.JiraBaseURL != "" {
		d.Register(notify.NewJiraBackend(cfg.JiraBaseURL, cfg.JiraProject, cfg.JiraUsername, cfg.JiraAPIToken))
	}

	return d
}
