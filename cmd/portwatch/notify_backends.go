package main

import (
	"log"

	"portwatch/internal/config"
	"portwatch/internal/notify"
)

// buildDispatcher constructs a Dispatcher wired with all configured backends.
func buildDispatcher(cfg *config.Config) *notify.Dispatcher {
	d := notify.New(log.Default())

	d.Add(notify.NewLogBackend(log.Default()))

	if cfg.WebhookURL != "" {
		d.Add(notify.NewWebhookBackend(cfg.WebhookURL))
	}

	if cfg.PagerDutyKey != "" {
		d.Add(notify.NewPagerDutyBackend(cfg.PagerDutyKey))
	}

	if cfg.SlackWebhookURL != "" {
		d.Add(notify.NewSlackBackend(cfg.SlackWebhookURL))
	}

	if cfg.EmailHost != "" && cfg.EmailFrom != "" && cfg.EmailTo != "" {
		d.Add(notify.NewEmailBackend(cfg.EmailHost, cfg.EmailFrom, cfg.EmailTo))
	}

	if cfg.SMSGatewayURL != "" && cfg.SMSFrom != "" && cfg.SMSTo != "" {
		d.Add(notify.NewSMSBackend(cfg.SMSGatewayURL, cfg.SMSAPIKey, cfg.SMSFrom, cfg.SMSTo))
	}

	return d
}
