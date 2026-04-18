package main

import (
	"github.com/user/portwatch/internal/notify"
)

func buildDispatcher(cfg interface{ GetNotify() map[string]string }) *notify.Dispatcher {
	settings := cfg.GetNotify()
	d := notify.New()

	d.Add(notify.NewLogBackend())

	if url, ok := settings["webhook_url"]; ok && url != "" {
		d.Add(notify.NewWebhookBackend(url))
	}
	if url, ok := settings["slack_url"]; ok && url != "" {
		d.Add(notify.NewSlackBackend(url))
	}
	if url, ok := settings["discord_url"]; ok && url != "" {
		d.Add(notify.NewDiscordBackend(url))
	}
	if url, ok := settings["teams_url"]; ok && url != "" {
		d.Add(notify.NewTeamsBackend(url))
	}
	if key, ok := settings["pagerduty_key"]; ok && key != "" {
		d.Add(notify.NewPagerDutyBackend(key))
	}
	if key, ok := settings["opsgenie_key"]; ok && key != "" {
		d.Add(notify.NewOpsGenieBackend(key))
	}
	if url, ok := settings["victorops_url"]; ok && url != "" {
		d.Add(notify.NewVictorOpsBackend(url))
	}
	if url, ok := settings["sms_url"]; ok && url != "" {
		d.Add(notify.NewSMSBackend(url))
	}
	if addr, ok := settings["email_addr"]; ok && addr != "" {
		smtp := settings["smtp_host"]
		d.Add(notify.NewEmailBackend(smtp, addr))
	}

	return d
}
