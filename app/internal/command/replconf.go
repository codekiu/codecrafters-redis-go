package command

import (
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/internal/protocol"
)

type ReplconfCommand struct{}

func NewReplconfCommand() *ReplconfCommand {
	return &ReplconfCommand{}
}

func (c *ReplconfCommand) Execute(conn net.Conn) error {
	_, err := conn.Write([]byte(protocol.T_SIMPLE_STRING + "OK" + protocol.CRLF))
	return err
}

func (c *ReplconfCommand) Name() string {
	return "replconf"
}

func (c *ReplconfCommand) ParseArguments(args []string) error {
	return nil
}
