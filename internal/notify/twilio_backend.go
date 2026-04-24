package notify

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/user/portwatch/internal/alert"
)

// TwilioBackend sends SMS alerts via the Twilio REST API.
type TwilioBackend struct {
	accountSID string
	authToken  string
	from       string
	to         string
	client     *http.Client
}

// NewTwilioBackend creates a TwilioBackend. accountSID, authToken, from and to
// are all required.
func NewTwilioBackend(accountSID, authToken, from, to string) *TwilioBackend {
	return &TwilioBackend{
		accountSID: accountSID,
		authToken:  authToken,
		from:       from,
		to:         to,
		client:     &http.Client{},
	}
}

func (t *TwilioBackend) Name() string { return "twilio" }

func (t *TwilioBackend) Send(event alert.Event) error {
	endpoint := fmt.Sprintf(
		"https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json",
		t.accountSID,
	)

	body := url.Values{}
	body.Set("From", t.from)
	body.Set("To", t.to)
	body.Set("Body", fmt.Sprintf("[portwatch] %s port %d", event.Type, event.Port))

	req, err := http.NewRequest(http.MethodPost, endpoint, strings.NewReader(body.Encode()))
	if err != nil {
		return fmt.Errorf("twilio: build request: %w", err)
	}
	req.SetBasicAuth(t.accountSID, t.authToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := t.client.Do(req)
	if err != nil {
		return fmt.Errorf("twilio: send: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("twilio: unexpected status %d", resp.StatusCode)
	}
	return nil
}
