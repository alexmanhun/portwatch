package notify

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// SMSBackend sends alerts via an HTTP SMS gateway (e.g. Twilio-compatible).
type SMSBackend struct {
	gatewayURL string
	apiKey     string
	from       string
	to         string
	client     *http.Client
}

// NewSMSBackend creates a new SMSBackend.
func NewSMSBackend(gatewayURL, apiKey, from, to string) *SMSBackend {
	return &SMSBackend{
		gatewayURL: gatewayURL,
		apiKey:     apiKey,
		from:       from,
		to:         to,
		client:     &http.Client{},
	}
}

func (s *SMSBackend) Name() string { return "sms" }

func (s *SMSBackend) Send(event Event) error {
	payload, err := json.Marshal(map[string]string{
		"from":    s.from,
		"to":      s.to,
		"message": fmt.Sprintf("[portwatch] %s: port %d", event.Type, event.Port),
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, s.gatewayURL, strings.NewReader(string(payload)))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.apiKey)

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("sms gateway returned status %d", resp.StatusCode)
	}
	return nil
}
