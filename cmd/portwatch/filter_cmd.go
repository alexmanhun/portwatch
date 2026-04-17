package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"portwatch/internal/history"
)

func runFilter(args []string) {
	fs := flag.NewFlagSet("filter", flag.ExitOnError)
	port := fs.Int("port", 0, "filter by port number")
	event := fs.String("event", "", "filter by event type (opened/closed)")
	sinceStr := fs.String("since", "", "filter records since duration ago (e.g. 24h)")
	untilStr := fs.String("until", "", "filter records until duration ago (e.g. 1h)")
	limit := fs.Int("limit", 0, "max number of records to return")
	format := fs.String("format", "text", "output format: text or csv")
	_ = fs.Parse(args)

	cfg := loadConfig()
	h, err := history.New(cfg.HistoryFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error opening history: %v\n", err)
		os.Exit(1)
	}

	opts := history.FilterOptions{
		Port:  *port,
		Event: *event,
		Limit: *limit,
	}
	if *sinceStr != "" {
		if d, err := time.ParseDuration(*sinceStr); err == nil {
			opts.Since = time.Now().Add(-d)
		}
	}
	if *untilStr != "" {
		if d, err := time.ParseDuration(*untilStr); err == nil {
			opts.Until = time.Now().Add(-d)
		}
	}

	records := history.Filter(h.All(), opts)

	switch *format {
	case "csv":
		history.ExportCSV(os.Stdout, records)
	default:
		history.ExportText(os.Stdout, records)
	}
}
