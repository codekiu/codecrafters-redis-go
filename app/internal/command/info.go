package command

import (
	"fmt"
	"net"
	"reflect"
	"strconv"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/internal/protocol"
	"github.com/codecrafters-io/redis-starter-go/app/internal/storage"
)

type InfoCommand struct {
	info    *storage.Information
	Content string
}

func NewInfoCommand(storage *storage.Information) *InfoCommand {
	return &InfoCommand{
		info: storage,
	}
}

func (c *InfoCommand) Execute(conn net.Conn) error {
	returnString := protocol.T_BULK_STRING
	var content string

	reflected := reflect.ValueOf(*c.info)
	typ := reflect.TypeOf(*c.info)

	for i := 0; i < reflected.NumField(); i++ {
		content += strings.ToLower(typ.Field(i).Name) + ":" + reflected.Field(i).String() + protocol.CRLF
	}

	returnString += strconv.Itoa(len(content)) + protocol.CRLF + content + protocol.CRLF

	_, err := conn.Write([]byte(returnString))
	return err
}

func (c *InfoCommand) Name() string {
	return "info"
}

func (c *InfoCommand) ParseArguments(args []string) error {
	fmt.Println(args)
	if len(args) < 5 {
		return fmt.Errorf("not enough args")
	}
	c.Content = args[4]

	return nil
}
