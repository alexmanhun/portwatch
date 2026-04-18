package notify

import (
	"fmt"
	"net/smtp"
	"strings"
)

// EmailConfig holds SMTP configuration for email notifications.
type EmailConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	To       []string
}

type emailBackend struct {
	cfg  EmailConfig
	send func(addr, from string, to []string, msg []byte) error
}

// NewEmailBackend creates a Backend that sends alerts via SMTP.
func NewEmailBackend(cfg EmailConfig) Backend {
	return &emailBackend{
		cfg:  cfg,
		send: smtp.SendMail,
	}
}

func (e *emailBackend) Name() string { return "email" }

func (e *emailBackend) Send(event string, detail string) error {
	addr := fmt.Sprintf("%s:%d", e.cfg.Host, e.cfg.Port)
	subject := fmt.Sprintf("Subject: [portwatch] %s\r\n", event)
	headers := strings.Join([]string{
		fmt.Sprintf("From: %s", e.cfg.From),
		fmt.Sprintf("To: %s", strings.Join(e.cfg.To, ",")),
		subject,
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=utf-8",
		"",
	}, "\r\n")
	body := headers + detail + "\r\n"

	var auth smtp.Auth
	if e.cfg.Username != "" {
		auth = smtp.PlainAuth("", e.cfg.Username, e.cfg.Password, e.cfg.Host)
	}
	return e.send(addr, e.cfg.From, e.cfg.To, []byte(body))
}
