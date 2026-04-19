package notify

import (
	"fmt"
	"log/syslog"
)

// SyslogBackend sends notifications to the local syslog daemon.
type SyslogBackend struct {
	tag      string
	priority syslog.Priority
	writer   *syslog.Writer
}

// NewSyslogBackend creates a SyslogBackend with the given tag and syslog priority.
func NewSyslogBackend(tag string, priority syslog.Priority) (*SyslogBackend, error) {
	w, err := syslog.New(priority, tag)
	if err != nil {
		return nil, fmt.Errorf("syslog: open: %w", err)
	}
	return &SyslogBackend{tag: tag, priority: priority, writer: w}, nil
}

func (s *SyslogBackend) Name() string { return "syslog" }

func (s *SyslogBackend) Send(event Event) error {
	msg := fmt.Sprintf("portwatch [%s] port %d: %s", event.Kind, event.Port, event.Message)
	if err := s.writer.Info(msg); err != nil {
		return fmt.Errorf("syslog: write: %w", err)
	}
	return nil
}

func (s *SyslogBackend) Close() error {
	if s.writer != nil {
		return s.writer.Close()
	}
	return nil
}
