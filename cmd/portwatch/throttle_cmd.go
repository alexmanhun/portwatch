package main

import (
	"fmt"
	"strconv"
	"time"

	"portwatch/internal/config"
	"portwatch/internal/notify"
)

// runThrottleTest exercises the throttle backend wrapper via the CLI.
// Usage: portwatch throttle-test <cooldown_seconds> <repeat_count>
func runThrottleTest(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: portwatch throttle-test <cooldown_seconds> <repeat_count>")
	}

	cooldownSec, err := strconv.Atoi(args[0])
	if err != nil || cooldownSec < 0 {
		return fmt.Errorf("invalid cooldown_seconds: %s", args[0])
	}

	repeat, err := strconv.Atoi(args[1])
	if err != nil || repeat < 1 {
		return fmt.Errorf("invalid repeat_count: %s", args[1])
	}

	cfg, _ := loadConfig()
	base := buildDispatcher(cfg)
	cooldown := time.Duration(cooldownSec) * time.Second
	tb := notify.NewThrottleBackend(base, cooldown)

	ev := testEvent(cfg)
	sent := 0
	for i := 0; i < repeat; i++ {
		if err := tb.Send(ev); err != nil {
			fmt.Printf("[throttle-test] error on attempt %d: %v\n", i+1, err)
			continue
		}
		sent++
	}

	fmt.Printf("[throttle-test] sent %d/%d events through throttle (cooldown=%ds)\n",
		sent, repeat, cooldownSec)
	return nil
}

func testEvent(cfg *config.Config) interface{ Type string } {
	// re-use the alert package helper indirectly via the existing notify path
	_ = cfg
	return nil
}
