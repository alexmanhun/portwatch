package main

import (
	"fmt"
	"os"

	"portwatch/internal/config"
)

func loadConfig() config.Config {
	cfg, err := config.Load("portwatch.json")
	if err != nil {
		cfg = config.Default()
	}
	if h := os.Getenv("PORTWATCH_HISTORY"); h != "" {
		cfg.HistoryFile = h
	}
	return cfg
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "report":
		runReport(os.Args[2:])
	case "filter":
		runFilter(os.Args[2:])
	case "run":
		runDaemon()
	case "help", "--help", "-h":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage: portwatch <command> [options]")
	fmt.Println("Commands:")
	fmt.Println("  run     Start the portwatch daemon")
	fmt.Println("  report  Display port history report")
	fmt.Println("  filter  Filter port history by criteria")
	fmt.Println("  help    Show this help message")
}

func runDaemon() {
	cfg := loadConfig()
	fmt.Printf("Starting portwatch daemon (ports %d-%d, interval %s)\n",
		cfg.StartPort, cfg.EndPort, cfg.Interval)
	select {}
}
