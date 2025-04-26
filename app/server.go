package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	// Uncomment this block to pass the first stage
	"net"
)

var dict = make(map[string]string)

type Command interface {
	Handle(conn net.Conn)
}

type pingCommand struct{}

func (c *pingCommand) Handle(conn net.Conn) {
	conn.Write([]byte(T_SIMPLE_STRING + "PONG" + CRLF))
}

type echoCommand struct {
	Content string
}

func (c *echoCommand) Handle(conn net.Conn) {
	conn.Write([]byte(T_SIMPLE_STRING + c.Content + CRLF))
}

type setCommand struct {
	Content []string
}

func (c *setCommand) Handle(conn net.Conn) {
	key := c.Content[4]
	value := c.Content[6]
	dict[key] = value
	if len(c.Content) > 10 && c.Content[8] == "px" {
		timeInInt, err := strconv.Atoi(c.Content[10])
		if err != nil {
			conn.Write([]byte(T_SIMPLE_STRING + "introduce a number to expire" + CRLF))
			return
		}
		go removeFromDict(key, time.Duration(timeInInt*int(time.Millisecond)))
	}
	conn.Write([]byte(T_SIMPLE_STRING + "OK" + CRLF))
}

type getCommand struct {
	Key string
}

func (c *getCommand) Handle(conn net.Conn) {
	if _, prs := dict[c.Key]; !prs {
		conn.Write([]byte(T_BULK_STRING + "-1" + CRLF))
		return
	}
	conn.Write([]byte(T_SIMPLE_STRING + dict[c.Key] + CRLF))
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	port := flag.Int("port", 6379, "Port for TCP server; 6379 as default")
	address := "0.0.0.0"

	fullAddress := address + ":" + strconv.Itoa(*port)

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
			if errors.Is(err, io.EOF) || strings.Contains(err.Error(), "closed pipe") {
				// Client closed connection normally or pipe was closed - this is expected in tests
				return
			}
			fmt.Println("Error reading from connection: ", err.Error())
			return
		}

		cmd, innerErr := parseCommand(string(buf[:n]))
		if innerErr != nil {
			conn.Write([]byte(T_SIMPLE_STRING + innerErr.Error() + CRLF))
			continue
		}

		cmd.Handle(conn)
	}
}

func parseCommand(request string) (Command, error) {
	messages := strings.Split(request, CRLF)
	arrayLength := messages[0]
	numElements, err := strconv.Atoi(arrayLength[1:])
	if err != nil {
		return nil, fmt.Errorf("wrong number of parameters'%v'", err)
	}

	cmd := strings.ToLower(messages[2])

	switch cmd {
	case "set":
		if numElements < 3 {
			return nil, errors.New("not enough parameters")
		}
		return &setCommand{Content: messages}, nil
	case "get":
		if numElements < 2 {
			return nil, errors.New("not enough parameters")
		}
		return &getCommand{Key: messages[4]}, nil
	case "echo":
		return &echoCommand{Content: messages[4]}, nil
	case "ping":
		return &pingCommand{}, nil
	}

	return nil, fmt.Errorf("unknown command '%s'", cmd)
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
	CRLF            = "\r\n"
)

func removeFromDict(key string, after time.Duration) {
	time.Sleep(after)
	delete(dict, key)
}
