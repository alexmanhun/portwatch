package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/alert"
)

// AzureServiceBusBackend sends notifications to an Azure Service Bus queue
// via the REST API using a SAS token.
type AzureServiceBusBackend struct {
	endpointURL string // full URL including queue path
	sasToken    string
	client      *http.Client
}

// NewAzureServiceBusBackend creates a new AzureServiceBusBackend.
// endpointURL should be of the form:
//
//	https://<namespace>.servicebus.windows.net/<queue>/messages
func NewAzureServiceBusBackend(endpointURL, sasToken string) *AzureServiceBusBackend {
	return &AzureServiceBusBackend{
		endpointURL: endpointURL,
		sasToken:    sasToken,
		client:      &http.Client{},
	}
}

func (a *AzureServiceBusBackend) Name() string { return "azureservicebus" }

func (a *AzureServiceBusBackend) Send(event alert.Event) error {
	payload := map[string]interface{}{
		"source": "portwatch",
		"type":   event.Type,
		"port":   event.Port,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, a.endpointURL, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", a.sasToken)

	resp, err := a.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("azureservicebus: unexpected status %d", resp.StatusCode)
	}
	return nil
}
