package main

import (
	"flag"
	"fmt"
	"os"

	"portwatch/internal/notify"
)

// runNotifyTest sends a test notification through all configured backends.
func runNotifyTest(args []string) {
	fs := flag.NewFlagSet("notify-test", flag.ExitOnError)
	webhook := fs.String("webhook", "", "Webhook URL to test")
	email := fs.String("email", "", "Email address to send test alert")
	smtpHost := fs.String("smtp-host", "localhost", "SMTP host")
	smtpPort := fs.Int("smtp-port", 25, "SMTP port")
	smtpUser := fs.String("smtp-user", "", "SMTP username")
	smtpPass := fs.String("smtp-pass", "", "SMTP password")
	smtpFrom := fs.String("smtp-from", "portwatch@localhost", "SMTP from address")
	_ = fs.Parse(args)

	dispatcher := notify.New()

	dispatcher.Register(notify.NewLogBackend(os.Stdout))

	if *webhook != "" {
		dispatcher.Register(notify.NewWebhookBackend(*webhook))
	}

	if *email != "" {
		cfg := notify.EmailConfig{
			Host:     *smtpHost,
			Port:     *smtpPort,
			Username: *smtpUser,
			Password: *smtpPass,
			From:     *smtpFrom,
			To:       []string{*email},
		}
		dispatcher.Register(notify.NewEmailBackend(cfg))
	}

	err := dispatcher.Dispatch("test_alert", "portwatch test notification")
	if err != nil {
		fmt.Fprintf(os.Stderr, "notify-test errors: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("notify-test: all backends dispatched successfully")
}
