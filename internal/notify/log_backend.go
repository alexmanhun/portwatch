package notify

import (
	"fmt"
	"io"
	"os"
	"time"
)

// LogBackend writes notifications as timestamped lines to a writer.
type LogBackend struct {
	out io.Writer
}

// NewLogBackend creates a LogBackend writing to w.
// If w is nil, os.Stdout is used.
func NewLogBackend(w io.Writer) *LogBackend {
	if w == nil {
		w = os.Stdout
	}
	return &LogBackend{out: w}
}

func (l *LogBackend) Name() string { return "log" }

func (l *LogBackend) Send(msg Message) error {
	_, err := fmt.Fprintf(
		l.out,
		"%s [%s] %s — %s (port %d)\n",
		time.Now().Format(time.RFC3339),
		msg.Level,
		msg.Title,
		msg.Body,
		msg.Port,
	)
	return err
}
