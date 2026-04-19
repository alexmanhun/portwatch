package notify

import (
	"fmt"
	"net"
	"time"
)

// SNMPBackend sends SNMP trap notifications.
type SNMPBackend struct {
	host      string
	port      int
	community string
}

// NewSNMPBackend creates a new SNMP trap backend.
func NewSNMPBackend(host string, port int, community string) *SNMPBackend {
	if port == 0 {
		port = 162
	}
	if community == "" {
		community = "public"
	}
	return &SNMPBackend{host: host, port: port, community: community}
}

func (s *SNMPBackend) Name() string {
	return "snmp"
}

// Send sends a minimal SNMP v1 trap UDP datagram with the alert message.
// This is a best-effort implementation suitable for simple alerting.
func (s *SNMPBackend) Send(event string, port int, meta string) error {
	addr := fmt.Sprintf("%s:%d", s.host, s.port)
	conn, err := net.DialTimeout("udp", addr, 3*time.Second)
	if err != nil {
		return fmt.Errorf("snmp: dial %s: %w", addr, err)
	}
	defer conn.Close()

	// Encode a minimal payload: community + event message as plain text.
	// Real SNMP traps use BER-encoded ASN.1; this stub sends a human-readable
	// datagram useful for testing and simple syslog-style receivers.
	payload := fmt.Sprintf("community=%s event=%s port=%d %s", s.community, event, port, meta)
	_, err = fmt.Fprint(conn, payload)
	if err != nil {
		return fmt.Errorf("snmp: send: %w", err)
	}
	return nil
}
