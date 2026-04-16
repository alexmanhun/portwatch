package monitor_test

import (
	"net"
	"strconv"
	"testing"
	"time"

	"portwatch/internal/monitor"
	"portwatch/internal/scanner"
)

func startServer(t *testing.T) (int, func()) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start server: %v", err)
	}
	port, _ := strconv.Atoi(ln.Addr().(*net.TCPAddr).Port.String())
	port = ln.Addr().(*net.TCPAddr).Port
	return port, func() { ln.Close() }
}

func TestMonitorDetectsNewPort(t *testing.T) {
	s := scanner.New(500 * time.Millisecond)

	var gotAdded, gotRemoved []int
	m := monitor.New(s, time.Second, func(added, removed []int) {
		gotAdded = added
		gotRemoved = removed
	})

	// First scan — no previous state, no callback expected.
	if err := m.Scan("127.0.0.1", 10000, 10100); err != nil {
		t.Fatalf("first scan error: %v", err)
	}

	// Open a port and scan again.
	port, close := startServer(t)
	defer close()

	if port < 10000 || port > 10100 {
		t.Skip("random port outside test range, skipping")
	}

	if err := m.Scan("127.0.0.1", 10000, 10100); err != nil {
		t.Fatalf("second scan error: %v", err)
	}

	if len(gotAdded) == 0 {
		t.Errorf("expected added ports, got none")
	}
	_ = gotRemoved
}

func TestMonitorNoChangeNoCallback(t *testing.T) {
	s := scanner.New(500 * time.Millisecond)
	called := false
	m := monitor.New(s, time.Second, func(added, removed []int) {
		called = true
	})

	for i := 0; i < 2; i++ {
		if err := m.Scan("127.0.0.1", 19900, 19910); err != nil {
			t.Fatalf("scan %d error: %v", i, err)
		}
	}

	if called {
		t.Error("onChange should not be called when ports are unchanged")
	}
}
