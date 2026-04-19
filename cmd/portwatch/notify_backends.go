package main

import (
	"portwatch/internal/config"
	"portwatch/internal/notify"
)

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
	if cfg.MQTTBrokerURL != "" {
		topic := cfg.MQTTTopic
		if topic == "" {
			topic = "portwatch/alerts"
		}
		d.Add(notify.NewMQTTBackend(cfg.MQTTBrokerURL, topic))
	}

	return d
}
