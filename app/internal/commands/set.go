package commands

import (
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/internal/protocol"
	"github.com/codecrafters-io/redis-starter-go/app/internal/storage"
)

type SetCommand struct {
	Storage *storage.MemoryStorage
	key     string
	value   string
	args    []string
}

func (c *SetCommand) Execute(conn net.Conn) error {
	if len(c.args) > 10 && c.args[8] == "px" {
		timeInInt, err := strconv.Atoi(c.args[10])
		if err != nil {
			return err
		}

		c.Storage.SetWithExpiry(c.key, c.value, time.Millisecond*time.Duration(timeInInt))
	}
	c.Storage.Set(c.key, c.value)

	conn.Write([]byte(protocol.T_SIMPLE_STRING + "OK" + protocol.CRLF))
	return nil
}

func (c *SetCommand) Name() string {
	return "set"
}

func (c *SetCommand) ParseArguments(args []string) error {
	fmt.Println("args", args)
	arrayLength := args[0]
	numElements, err := strconv.Atoi(arrayLength[1:])
	if err != nil {
		return fmt.Errorf("wrong number of parameters'%v'", err)
	}

	if numElements < 3 {
		return fmt.Errorf("not enough args")
	}

	c.key = args[4]
	c.value = args[6]
	c.args = args

	return nil
}
