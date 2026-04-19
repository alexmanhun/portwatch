package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// AMQPBackend sends alerts to a RabbitMQ management HTTP API exchange.
type AMQPBackend struct {
	baseURL    string
	vhost      string
	exchange   string
	routingKey string
	username   string
	password   string
	client     *http.Client
}

// NewAMQPBackend creates a new AMQPBackend targeting the RabbitMQ management API.
func NewAMQPBackend(baseURL, vhost, exchange, routingKey, username, password string) *AMQPBackend {
	return &AMQPBackend{
		baseURL:    baseURL,
		vhost:      vhost,
		exchange:   exchange,
		routingKey: routingKey,
		username:   username,
		password:   password,
		client:     &http.Client{},
	}
}

func (a *AMQPBackend) Name() string { return "amqp" }

func (a *AMQPBackend) Send(event Event) error {
	payload := map[string]interface{}{
		"properties":       map[string]interface{}{"content_type": "application/json"},
		"routing_key":      a.routingKey,
		"payload":          fmt.Sprintf("%s port %d", event.Type, event.Port),
		"payload_encoding": "string",
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%s/api/exchanges/%s/%s/publish", a.baseURL, a.vhost, a.exchange)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.SetBasicAuth(a.username, a.password)
	req.Header.Set("Content-Type", "application/json")
	resp, err := a.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("amqp backend: unexpected status %d", resp.StatusCode)
	}
	return nil
}
