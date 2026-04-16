package history

import (
	"os"
	"path/filepath"
	"testing"
)

func tempFile(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "history.json")
}

func TestRecordAndRetrieve(t *testing.T) {
	h, err := New(tempFile(t))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := h.Record(8080, "opened"); err != nil {
		t.Fatalf("Record: %v", err)
	}
	if err := h.Record(8080, "closed"); err != nil {
		t.Fatalf("Record: %v", err)
	}
	entries := h.Entries()
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Port != 8080 || entries[0].Event != "opened" {
		t.Errorf("unexpected first entry: %+v", entries[0])
	}
}

func TestPersistence(t *testing.T) {
	path := tempFile(t)
	h, _ := New(path)
	h.Record(9090, "opened")

	h2, err := New(path)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	entries := h2.Entries()
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry after reload, got %d", len(entries))
	}
	if entries[0].Port != 9090 {
		t.Errorf("unexpected port: %d", entries[0].Port)
	}
}

func TestMissingFileIsOK(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nonexistent.json")
	h, err := New(path)
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if len(h.Entries()) != 0 {
		t.Error("expected empty entries")
	}
}

func TestTimestampSet(t *testing.T) {
	h, _ := New(tempFile(t))
	h.Record(443, "opened")
	e := h.Entries()[0]
	if e.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestNoFileCreatedUntilRecord(t *testing.T) {
	path := filepath.Join(t.TempDir(), "noop.json")
	New(path)
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Error("file should not exist before any Record call")
	}
}
