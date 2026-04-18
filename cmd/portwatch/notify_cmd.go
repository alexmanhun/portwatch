package main

import (
	"flag"
	"fmt"
	"os"

	"portwatch/internal/notify"
)

func runNotifyTest(args []string) {
	fs := flag.NewFlagSet("notify-test", flag.ExitOnError)
	webhook := fs.String("webhook", "", "Webhook URL to test")
	slack := fs.String("slack", "", "Slack webhook URL to test")
	email := fs.String("email", "", "Email address to test (requires SMTP config)")
	fs.Parse(args)

	dispatcher := notify.New()

	if *webhook != "" {
		dispatcher.Register(notify.NewWebhookBackend(*webhook))
		fmt.Println("Registered webhook backend")
	}
	if *slack != "" {
		dispatcher.Register(notify.NewSlackBackend(*slack))
		fmt.Println("Registered Slack backend")
	}
	if *email != "" {
		cfg := loadConfig()
		be := notify.NewEmailBackend(
			cfg.SMTPHost, cfg.SMTPPort,
			cfg.SMTPUser, cfg.SMTPPass,
			cfg.AlertFrom, *email,
		)
		dispatcher.Register(be)
		fmt.Println("Registered email backend")
	}

	dispatcher.Register(notify.NewLogBackend(os.Stdout))

	event := notify.Event{
		Type:    "test",
		Port:    0,
		Message: "portwatch notify test — if you see this, notifications are working",
	}

	if err := dispatcher.Dispatch(event); err != nil {
		fmt.Fprintf(os.Stderr, "notify-test errors: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Notify test complete.")
}
