package history

import (
	"encoding/csv"
	"fmt"
	"io"
	"time"
)

// ExportCSV writes history records to w in CSV format.
// Columns: timestamp, port, event
func ExportCSV(records []Record, w io.Writer) error {
	cw := csv.NewWriter(w)
	if err := cw.Write([]string{"timestamp", "port", "event"}); err != nil {
		return fmt.Errorf("export csv header: %w", err)
	}
	for _, r := range records {
		row := []string{
			r.Timestamp.UTC().Format(time.RFC3339),
			fmt.Sprintf("%d", r.Port),
			string(r.Event),
		}
		if err := cw.Write(row); err != nil {
			return fmt.Errorf("export csv row: %w", err)
		}
	}
	cw.Flush()
	return cw.Error()
}

// ExportText writes history records as human-readable lines to w.
func ExportText(records []Record, w io.Writer) error {
	for _, r := range records {
		line := fmt.Sprintf("%s  port=%-6d  event=%s\n",
			r.Timestamp.UTC().Format(time.RFC3339),
			r.Port,
			r.Event,
		)
		if _, err := io.WriteString(w, line); err != nil {
			return fmt.Errorf("export text: %w", err)
		}
	}
	return nil
}
