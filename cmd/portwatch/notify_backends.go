package main

import (
	"os"

	"portwatch/internal/notify"
)

// buildDispatcher constructs a Dispatcher from environment-configured backends.
// Supported backends are enabled by setting the corresponding environment variables.
func buildDispatcher() *notify.Dispatcher {
	d := notify.New()

	if webhookURL := os.Getenv("PORTWATCH_WEBHOOK_URL"); webhookURL != "" {
		d.Register(notify.NewWebhookBackend(webhookURL))
	}

	if slackURL := os.Getenv("PORTWATCH_SLACK_URL"); slackURL != "" {
		d.Register(notify.NewSlackBackend(slackURL))
	}

	if pdKey := os.Getenv("PORTWATCH_PAGERDUTY_KEY"); pdKey != "" {
		d.Register(notify.NewPagerDutyBackend(pdKey))
	}

	if smtpHost := os.Getenv("PORTWATCH_SMTP_HOST"); smtpHost != "" {
		to := os.Getenv("PORTWATCH_ALERT_EMAIL")
		from := os.Getenv("PORTWATCH_FROM_EMAIL")
		if to != "" && from != "" {
			d.Register(notify.NewEmailBackend(smtpHost, from, to))
		}
	}

	// Always attach log backend as fallback.
	d.Register(notify.NewLogBackend(os.Stderr))

	return d
}
