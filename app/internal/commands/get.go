package commands

import (
	"fmt"
	"net"
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/app/internal/protocol"
	"github.com/codecrafters-io/redis-starter-go/app/internal/storage"
)

type GetCommand struct {
	Storage *storage.MemoryStorage
	key     string
}

func (c *GetCommand) Execute(conn net.Conn) error {
	content, exists := c.Storage.Get(c.key)
	if !exists {
		conn.Write([]byte(protocol.T_BULK_STRING + "-1" + protocol.CRLF))
		return nil
	}
	_, err := conn.Write([]byte(protocol.T_SIMPLE_STRING + content + protocol.CRLF))
	return err
}

func (c *GetCommand) Name() string {
	return "get"
}

func (c *GetCommand) ParseArguments(args []string) error {
	arrayLength := args[0]
	numElements, err := strconv.Atoi(arrayLength[1:])
	if err != nil {
		return fmt.Errorf("wrong number of parameters'%v'", err)
	}
	if numElements < 2 {
		return fmt.Errorf("not enough args")
	}

	c.key = args[4]

	return nil
}
