package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"portwatch/internal/alert"
	"portwatch/internal/config"
	"portwatch/internal/monitor"
	"portwatch/internal/scanner"
)

func main() {
	configPath := flag.String("config", "", "path to JSON config file")
	flag.Parse()

	var cfg *config.Config
	var err error

	if *configPath != "" {
		cfg, err = config.Load(*configPath)
		if err != nil {
			log.Fatalf("failed to load config: %v", err)
		}
	} else {
		cfg = config.Default()
	}

	var alertOut *os.File
	if cfg.AlertOutput != "" {
		alertOut, err = os.OpenFile(cfg.AlertOutput, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("failed to open		defer alertOut.Close()
	}

	notifier := alert.New(alertOut)
	sc := scanner.New(cfg.PortRange.Start, cfg.PortRange.End)

	mon := monitor.New(sc, cfg.Interval, func(opened, closed []int) {
		for _, p := range opened {
			notifier.Notify(alert.NewPortEvent(p))
		}
		for _, p := range closed {
			notifier.Notify(alert.ClosedPortEvent(p))
		}
	})

	fmt.Printf("portwatch starting — scanning ports %d-%d every %v\n",
		cfg.PortRange.Start, cfg.PortRange.End, cfg.Interval)

	mon.Start()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	fmt.Println("\nportwatch stopping")
	mon.Stop()
}
