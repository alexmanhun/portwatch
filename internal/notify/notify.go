// Package notify provides pluggable notification backends for portwatch.
package notify

import (
	"fmt"
	"io"
	"os"
)

// Level represents the severity of a notification.
type Level string

const (
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelAlert Level = "ALERT"
)

// Message holds the data for a single notification.
type Message struct {
	Level   Level
	Title   string
	Body    string
	Port    int
	Event   string
}

// Backend is the interface that all notification backends must implement.
type Backend interface {
	Send(msg Message) error
	Name() string
}

// Dispatcher fans out a message to all registered backends.
type Dispatcher struct {
	backends []Backend
	errOut   io.Writer
}

// New creates a Dispatcher with the given backends.
func New(backends ...Backend) *Dispatcher {
	return &Dispatcher{backends: backends, errOut: os.Stderr}
}

// SetErrOut redirects error output (useful in tests).
func (d *Dispatcher) SetErrOut(w io.Writer) {
	d.errOut = w
}

// Dispatch sends msg to every registered backend, logging errors.
func (d *Dispatcher) Dispatch(msg Message) {
	for _, b := range d.backends {
		if err := b.Send(msg); err != nil {
			fmt.Fprintf(d.errOut, "[notify] backend %s error: %v\n", b.Name(), err)
		}
	}
}
