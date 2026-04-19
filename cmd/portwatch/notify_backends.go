package main

import (
	"log/syslog"
	"os"

	"portwatch/internal/config"
	"portwatch/internal/notify"
)

func buildDispatcher(cfg *config.Config) *notify.Dispatcher {
	d := notify.New()

	// Always include the log backend.
	d.Add(notify.NewLogBackend(os.Stderr))

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
	if cfg.SyslogTag != "" {
		pri := syslog.LOG_INFO | syslog.LOG_DAEMON
		if b, err := notify.NewSyslogBackend(cfg.SyslogTag, pri); err == nil {
			d.Add(b)
		}
	}

	return d
}
