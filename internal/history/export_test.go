package history

import (
	"bytes"
	"encoding/csv"
	"strings"
	"testing"
	"time"
)

func sampleRecords() []Record {
	base := time, 15, 10, 0, 0, 0, time.UTC)
	return []Record{
		{Port: 8080, Event: "new",
		{Port: 443, Event: "closed", Timestamp: base.Add(time.Minute)},
	}
}

func TestExportCSVHeader(t *testing.T) {
	var buf bytes.Buffer
	if err := ExportCSV(sampleRecords(), &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	r := csv.NewReader(&buf)
	rows, err := r.ReadAll()
	if err != nil {
		t.Fatalf("invalid csv: %v", err)
	}
	if len(rows) != 3 {
		t.Fatalf("expected 3 rows (header+2), got %d", len(rows))
	}
	if rows[0][0] != "timestamp" || rows[0][1] != "port" || rows[0][2] != "event" {
		t.Errorf("unexpected header: %v", rows[0])
	}
}

func TestExportCSVValues(t *testing.T) {
	var buf bytes.Buffer
	_ = ExportCSV(sampleRecords(), &buf)
	r := csv.NewReader(&buf)
	rows, _ := r.ReadAll()
	if rows[1][1] != "8080" {
		t.Errorf("expected port 8080, got %s", rows[1][1])
	}
	if rows[2][2] != "closed" {
		t.Errorf("expected event closed, got %s", rows[2][2])
	}
}

func TestExportCSVEmpty(t *testing.T) {
	var buf bytes.Buffer
	if err := ExportCSV(nil, &buf); err != nil {
		t.Fatalf("unexpected error on empty: %v", err)
	}
}

func TestExportTextFormat(t *testing.T) {
	var buf bytes.Buffer
	if err := ExportText(sampleRecords(), &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "port=8080") {
		t.Errorf("expected port=8080 in output: %s", out)
	}
	if !strings.Contains(out, "event=new") {
		t.Errorf("expected event=new in output: %s", out)
	}
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 2 {
		t.Errorf("expected 2 lines, got %d", len(lines))
	}
}
