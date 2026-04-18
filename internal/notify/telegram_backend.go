package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// TelegramBackend sends notifications via the Telegram Bot API.
type TelegramBackend struct {
	botToken string
	chatID   string
	apiBase  string
}

// NewTelegramBackend creates a new TelegramBackend.
func NewTelegramBackend(botToken, chatID string) *TelegramBackend {
	return &TelegramBackend{
		botToken: botToken,
		chatID:   chatID,
		apiBase:  "https://api.telegram.org",
	}
}

func (t *TelegramBackend) Name() string { return "telegram" }

func (t *TelegramBackend) Send(event Event) error {
	url := fmt.Sprintf("%s/bot%s/sendMessage", t.apiBase, t.botToken)

	payload := map[string]string{
		"chat_id": t.chatID,
		"text":    fmt.Sprintf("[portwatch] %s: port %d", event.Type, event.Port),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("telegram: marshal: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("telegram: request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("telegram: unexpected status %d", resp.StatusCode)
	}
	return nil
}
