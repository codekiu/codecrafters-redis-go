package main

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	// Uncomment this block to pass the first stage
	"net"
)

type Command interface {
	Handle(conn net.Conn)
}

type pingCommand struct{}

func (c *pingCommand) Handle(conn net.Conn) {
	conn.Write([]byte("+PONG\r\n"))
}

type echoCommand struct {
	Content string
}

func (c *echoCommand) Handle(conn net.Conn) {
	conn.Write([]byte("+" + c.Content + "\r\n"))
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	port := 6379
	address := "0.0.0.0"

	fullAddress := address + ":" + strconv.Itoa(port)

	listener, err := net.Listen("tcp", fullAddress)
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		return
	}

	defer listener.Close()

	fmt.Printf("Server is listening in port: %d\n", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			continue
		}

		// Handle client connection
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	// Read data
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return // Client closed the connection normally
			}
			fmt.Println("Error reading from connection: ", err.Error())
			return
		}

		cmd, innerErr := parseCommand(string(buf[:n]))
		if innerErr != nil {
			fmt.Println("Error parsing command: ", innerErr.Error())
			return
		}

		cmd.Handle(conn)
	}
}

func parseCommand(request string) (Command, error) {
	messages := strings.Split(request, CRLF)
	fmt.Println("Messages: ", messages)
	cmd := strings.ToLower(messages[2])

	fmt.Println("CMD: ", cmd)
	switch cmd {
	case "echo":
		return &echoCommand{Content: messages[4]}, nil
	case "ping":
		return &pingCommand{}, nil
	}

	return nil, errors.New("no command to parse")
}

const (
	T_SIMPLE_STRING = "+"
	T_SIMPLE_ERROR  = "-"
	T_INTEGER       = ":"
	T_BULK_STRING   = "$"
	T_ARRAY         = "*"
	T_NULL          = "_"
	T_BOOLEAN       = "#"
	T_MAP           = "%"
)

const CRLF = "\r\n"
