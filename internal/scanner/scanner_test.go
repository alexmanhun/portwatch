package scanner

import (
	"net"
	"strconv"
	"testing"
)

// startTestServer opens a TCP listener on a random port and returns it.
func startTestServer(t *testing.T) (*net.TCPListener, int) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start test server: %v", err)
	}
	_, portStr, _ := net.SplitHostPort(ln.Addr().String())
	port, _ := strconv.Atoi(portStr)
	return ln.(*net.TCPListener), port
}

func TestScanDetectsOpenPort(t *testing.T) {
	ln, port := startTestServer(t)
	defer ln.Close()

	s := New("127.0.0.1")
	ports, err := s.Scan(port, port)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ports) != 1 || ports[0].Number != port {
		t.Errorf("expected port %d to be open, got %v", port, ports)
	}
}

func TestScanInvalidRange(t *testing.T) {
	s := New("127.0.0.1")
	_, err := s.Scan(9000, 8000)
	if err == nil {
		t.Error("expected error for invalid range, got nil")
	}
}

func TestScanClosedPort(t *testing.T) {
	// Port 1 is almost certainly closed in test environments.
	s := New("127.0.0.1")
	ports, err := s.Scan(1, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ports) != 0 {
		t.Errorf("expected no open ports, got %v", ports)
	}
}
