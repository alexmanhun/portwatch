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
	if cfg.SlackWebhook != "" {
		d.Add(notify.NewSlackBackend(cfg.SlackWebhook))
	}
	if cfg.KafkaProxyURL != "" && cfg.KafkaTopic != "" {
		d.Add(notify.NewKafkaBackend(cfg.KafkaProxyURL, cfg.KafkaTopic))
	}

	return d
}
