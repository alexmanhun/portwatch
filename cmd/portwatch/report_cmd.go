package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"portwatch/internal/history"
)

// runReport loads history from path and prints a summary report.
// Flags: --since duration, --port int (0 = all), --last n
func runReport(historyPath string, args []string) error {
	fs := flag.NewFlagSet("report", flag.ContinueOnError)
	sinceStr := fs.String("since", "", "show events since duration ago, e.g. 24h")
	port := fs.Int("port", 0, "filter by port (0 = all)")
	last := fs.Int("last", 0, "show only last N records")

	if err := fs.Parse(args); err != nil {
		return err
	}

	h, err := history.New(historyPath)
	if err != nil {
		return fmt.Errorf("open history: %w", err)
	}

	q := h.Query()

	if *port != 0 {
		q = q.ByPort(*port)
	}
	if *sinceStr != "" {
		d, err := time.ParseDuration(*sinceStr)
		if err != nil {
			return fmt.Errorf("invalid duration %q: %w", *sinceStr, err)
		}
		q = q.Since(time.Now().Add(-d))
	}
	if *last > 0 {
		q = q.Last(*last)
	}

	records := q.All()
	summary := history.Summarize(records)
	history.PrintReport(os.Stdout, summary)
	return nil
}
