package command

import (
	"fmt"
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/internal/protocol"
)

type EchoCommand struct {
	Content string
}

func NewEchoCommand() *EchoCommand {
	return &EchoCommand{}
}

func (c *EchoCommand) Execute(conn net.Conn) error {
	_, err := conn.Write([]byte(protocol.T_SIMPLE_STRING + c.Content + protocol.CRLF))
	return err
}

func (c *EchoCommand) Name() string {
	return "echo"
}

func (c *EchoCommand) ParseArguments(args []string) error {
	if len(args) < 4 {
		return fmt.Errorf("not enough args")
	}
	c.Content = args[4]

	return nil
}
