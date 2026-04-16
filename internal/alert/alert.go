package alert

import (
	"fmt"
	"io"
	"os"
	"time"
)

// Level represents the severity of an alert.
type Level string

const (
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelAlert Level = "ALERT"
)

// Event describes a port change event.
type Event struct {
	Timestamp time.Time
	Level     Level
	Port      int
	Message   string
}

// Notifier sends alerts for port change events.
type Notifier struct {
	out io.Writer
}

// New creates a Notifier that writes to the given writer.
// If w is nil, os.Stdout is used.
func New(w io.Writer) *Notifier {
	if w == nil {
		w = os.Stdout
	}
	return &Notifier{out: w}
}

// Notify formats and writes an alert event.
func (n *Notifier) Notify(e Event) error {
	_, err := fmt.Fprintf(
		n.out,
		"[%s] %s port=%d msg=%q\n",
		e.Timestamp.Format(time.RFC3339),
		e.Level,
		e.Port,
		e.Message,
	)
	return err
}

// NewPortEvent returns an ALERT Event for a newly opened port.
func NewPortEvent(port int) Event {
	return Event{
		Timestamp: time.Now(),
		Level:     LevelAlert,
		Port:      port,
		Message:   fmt.Sprintf("new open port detected: %d", port),
	}
}

// ClosedPortEvent returns a WARN Event for a port that has closed.
func ClosedPortEvent(port int) Event {
	return Event{
		Timestamp: time.Now(),
		Level:     LevelWarn,
		Port:      port,
		Message:   fmt.Sprintf("port closed: %d", port),
	}
}
