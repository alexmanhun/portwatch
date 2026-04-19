package notify

import (
	"net"
	"strings"
	"testing"
	"time"
)

func listenUDP(t *testing.T) (*net.UDPConn, int) {
	t.Helper()
	addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		t.Fatal(err)
	}
	return conn, conn.LocalAddr().(*net.UDPAddr).Port
}

func TestSNMPBackendName(t *testing.T) {
	b := NewSNMPBackend("localhost", 162, "public")
	if b.Name() != "snmp" {
		t.Errorf("expected snmp, got %s", b.Name())
	}
}

func TestSNMPBackendSendsPayload(t *testing.T) {
	conn, port := listenUDP(t)
	defer conn.Close()

	b := NewSNMPBackend("127.0.0.1", port, "testcommunity")
	err := b.Send("new_port", 8080, "proto=tcp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	buf := make([]byte, 512)
	n, _, err := conn.ReadFromUDP(buf)
	if err != nil {
		t.Fatalf("read udp: %v", err)
	}
	got := string(buf[:n])
	if !strings.Contains(got, "testcommunity") {
		t.Errorf("expected community in payload, got: %s", got)
	}
	if !strings.Contains(got, "new_port") {
		t.Errorf("expected event in payload, got: %s", got)
	}
	if !strings.Contains(got, "8080") {
		t.Errorf("expected port in payload, got: %s", got)
	}
}

func TestSNMPBackendDefaultPort(t *testing.T) {
	b := NewSNMPBackend("localhost", 0, "")
	if b.port != 162 {
		t.Errorf("expected default port 162, got %d", b.port)
	}
	if b.community != "public" {
		t.Errorf("expected default community public, got %s", b.community)
	}
}

func TestSNMPBackendBadHost(t *testing.T) {
	b := NewSNMPBackend("256.256.256.256", 162, "public")
	err := b.Send("new_port", 22, "")
	if err == nil {
		t.Error("expected error for bad host")
	}
}
