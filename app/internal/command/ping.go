package command

import (
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/internal/protocol"
)

type PingCommand struct{}

func NewPingCommand() *PingCommand {
	return &PingCommand{}
}

func (c *PingCommand) Execute(conn net.Conn) error {
	_, err := conn.Write([]byte(protocol.T_SIMPLE_STRING + "PONG" + protocol.CRLF))
	return err
}

func (c *PingCommand) Name() string {
	return "ping"
}

func (c *PingCommand) ParseArguments(args []string) error {
	return nil
}
