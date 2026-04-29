package main

import (
	"testing"
)

func TestRunThrottleTestMissingArgs(t *testing.T) {
	err := runThrottleTest([]string{})
	if err == nil {
		t.Fatal("expected error for missing args")
	}
}

func TestRunThrottleTestInvalidCooldown(t *testing.T) {
	err := runThrottleTest([]string{"abc", "3"})
	if err == nil {
		t.Fatal("expected error for non-numeric cooldown")
	}
}

func TestRunThrottleTestInvalidRepeat(t *testing.T) {
	err := runThrottleTest([]string{"1", "xyz"})
	if err == nil {
		t.Fatal("expected error for non-numeric repeat")
	}
}

func TestRunThrottleTestZeroRepeat(t *testing.T) {
	err := runThrottleTest([]string{"0", "0"})
	if err == nil {
		t.Fatal("expected error for zero repeat count")
	}
}

func TestRunThrottleTestNegativeCooldown(t *testing.T) {
	err := runThrottleTest([]string{"-1", "1"})
	if err == nil {
		t.Fatal("expected error for negative cooldown")
	}
}
