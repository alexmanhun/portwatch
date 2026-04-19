package notify

import (
	"fmt"
	"net"
	"time"
)

// IRCBackend sends notifications to an IRC channel via a simple TCP connection.
type IRCBackend struct {
	server  string
	nick    string
	channel string
}

// NewIRCBackend creates a new IRCBackend.
// server should be in host:port format, e.g. "irc.libera.chat:6667".
func NewIRCBackend(server, nick, channel string) *IRCBackend {
	return &IRCBackend{server: server, nick: nick, channel: channel}
}

func (b *IRCBackend) Name() string { return "irc" }

func (b *IRCBackend) Send(event string, port int) error {
	conn, err := net.DialTimeout("tcp", b.server, 5*time.Second)
	if err != nil {
		return fmt.Errorf("irc: dial %s: %w", b.server, err)
	}
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(10 * time.Second))

	send := func(msg string) error {
		_, err := fmt.Fprintf(conn, "%s\r\n", msg)
		return err
	}

	if err := send(fmt.Sprintf("NICK %s", b.nick)); err != nil {
		return err
	}
	if err := send(fmt.Sprintf("USER %s 0 * :portwatch", b.nick)); err != nil {
		return err
	}
	if err := send(fmt.Sprintf("JOIN %s", b.channel)); err != nil {
		return err
	}
	msg := fmt.Sprintf("[portwatch] %s on port %d", event, port)
	if err := send(fmt.Sprintf("PRIVMSG %s :%s", b.channel, msg)); err != nil {
		return err
	}
	return send("QUIT")
}
