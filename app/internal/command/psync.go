package command

import (
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/internal/protocol"
)

type PsyncCommand struct{}

func NewPsyncCommand() *PsyncCommand {
	return &PsyncCommand{}
}

func (c *PsyncCommand) Execute(conn net.Conn) error {
	_, err := conn.Write([]byte(protocol.T_SIMPLE_STRING + "FULLRESYNC 1 0" + protocol.CRLF))
	return err
}

func (c *PsyncCommand) Name() string {
	return "psync"
}

func (c *PsyncCommand) ParseArguments(args []string) error {
	return nil
}
