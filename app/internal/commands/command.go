package commands

import "net"

type Command interface {
	Execute(conn net.Conn) error

	Name() string

	ParseArguments(args []string) error
}
