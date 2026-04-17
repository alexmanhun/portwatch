package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"portwatch/internal/history"
)

func setupFilterHistory(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "history.json")
	h, err := history.New(p)
	if err != nil {
		t.Fatal(err)
	}
	h.Record(80, "opened")
	h.Record(443, "opened")
	h.Record(80, "closed")
	return p
}

func TestFilterCmdTextOutput(t *testing.T) {
	p := setupFilterHistory(t)
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	t.Setenv("PORTWATCH_HISTORY", p)
	runFilter([]string{"--port=80", "--format=text"})
	w.Close()
	os.Stdout = old
	buf := make([]byte, 1024)
	n, _ := r.Read(buf)
	out := string(buf[:n])
	if out == "" {
		t.Error("expected output, got empty string")
	}
}

func TestFilterCmdCSVOutput(t *testing.T) {
	p := setupFilterHistory(t)
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	t.Setenv("PORTWATCH_HISTORY", p)
	runFilter([]string{"--format=csv"})
	w.Close()
	os.Stdout = old
	buf := make([]byte, 2048)
	n, _ := r.Read(buf)
	out := string(buf[:n])
	if len(out) == 0 {
		t.Error("expected CSV output")
	}
}

func TestFilterCmdSince(t *testing.T) {
	_ = time.Now() // ensure time package used
	p := setupFilterHistory(t)
	t.Setenv("PORTWATCH_HISTORY", p)
	// should not panic
	runFilter([]string{"--since=999h", "--format=text"})
}
