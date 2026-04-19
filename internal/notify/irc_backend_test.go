package notify

import (
	"bufio"
	"net"
	"strings"
	"testing"
)

func TestIRCBackendName(t *testing.T) {
	b := NewIRCBackend("localhost:6667", "portwatch", "#alerts")
	if b.Name() != "irc" {
		t.Fatalf("expected irc, got %s", b.Name())
	}
}

func TestIRCBackendSendsMessage(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()

	received := make(chan []string, 1)
	go func() {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		defer conn.Close()
		var lines []string
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
		received <- lines
	}()

	b := NewIRCBackend(ln.Addr().String(), "portwatch", "#alerts")
	if err := b.Send("new_port", 8080); err != nil {
		t.Fatalf("Send returned error: %v", err)
	}

	lines := <-received
	full := strings.Join(lines, "\n")
	if !strings.Contains(full, "PRIVMSG #alerts") {
		t.Errorf("expected PRIVMSG, got: %s", full)
	}
	if !strings.Contains(full, "new_port") || !strings.Contains(full, "8080") {
		t.Errorf("message missing event/port: %s", full)
	}
}

func TestIRCBackendBadURL(t *testing.T) {
	b := NewIRCBackend("127.0.0.1:1", "portwatch", "#alerts")
	if err := b.Send("new_port", 22); err == nil {
		t.Fatal("expected error for unreachable server")
	}
}
