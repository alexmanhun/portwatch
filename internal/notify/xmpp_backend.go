package notify

import (
	"fmt"
	"net/http"
	"strings"
)

// XMPPBackend sends alerts via an XMPP HTTP gateway.
type XMPPBackend struct {
	gatewayURL string
	to         string
	client     *http.Client
}

// NewXMPPBackend creates a new XMPPBackend.
// gatewayURL should be an HTTP endpoint that accepts POST with a `to` and `body` form field.
func NewXMPPBackend(gatewayURL, to string) *XMPPBackend {
	return &XMPPBackend{
		gatewayURL: gatewayURL,
		to:         to,
		client:     &http.Client{},
	}
}

func (x *XMPPBackend) Name() string { return "xmpp" }

func (x *XMPPBackend) Send(event Event) error {
	body := fmt.Sprintf("[portwatch] %s port %d", event.Type, event.Port)
	form := strings.NewReader(fmt.Sprintf("to=%s&body=%s", x.to, body))
	resp, err := x.client.Post(x.gatewayURL, "application/x-www-form-urlencoded", form)
	if err != nil {
		return fmt.Errorf("xmpp: request failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("xmpp: unexpected status %d", resp.StatusCode)
	}
	return nil
}
