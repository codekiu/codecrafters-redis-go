package command

import (
	"encoding/hex"
	"net"
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/app/internal/protocol"
)

type PsyncCommand struct{}

func NewPsyncCommand() *PsyncCommand {
	return &PsyncCommand{}
}

func (c *PsyncCommand) Execute(conn net.Conn) error {
	_, err := conn.Write([]byte(protocol.T_SIMPLE_STRING + "FULLRESYNC 1 0" + protocol.CRLF))
	go sendRdbFile(conn)
	return err
}

func (c *PsyncCommand) Name() string {
	return "psync"
}

func (c *PsyncCommand) ParseArguments(args []string) error {
	return nil
}

func sendRdbFile(conn net.Conn) {
	rdbFile, _ := hex.DecodeString("524544495330303131fa0972656469732d76657205372e322e30fa0a72656469732d62697473c040fa056374696d65c26d08bc65fa08757365642d6d656dc2b0c41000fa08616f662d62617365c000fff06e3bfec0ff5aa2")

	conn.Write([]byte(protocol.T_BULK_STRING + strconv.Itoa(len(rdbFile)) + protocol.CRLF + string(rdbFile)))
}
