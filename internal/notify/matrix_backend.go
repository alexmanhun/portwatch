package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// MatrixBackend sends notifications to a Matrix room via the Client-Server API.
type MatrixBackend struct {
	homeserver string
	token      string
	roomID     string
	client     *http.Client
}

// NewMatrixBackend creates a new MatrixBackend.
// homeserver is the base URL (e.g. https://matrix.org),
// token is the access token, roomID is the Matrix room ID.
func NewMatrixBackend(homeserver, token, roomID string) *MatrixBackend {
	return &MatrixBackend{
		homeserver: homeserver,
		token:      token,
		roomID:     roomID,
		client:     &http.Client{},
	}
}

func (m *MatrixBackend) Name() string { return "matrix" }

func (m *MatrixBackend) Send(event Event) error {
	url := fmt.Sprintf("%s/_matrix/client/v3/rooms/%s/send/m.room.message",
		m.homeserver, m.roomID)

	body, err := json.Marshal(map[string]string{
		"msgtype": "m.text",
		"body":    fmt.Sprintf("[portwatch] %s: port %d", event.Type, event.Port),
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+m.token)

	resp, err := m.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("matrix: unexpected status %d", resp.StatusCode)
	}
	return nil
}
