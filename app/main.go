package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/internal/commands"
	"github.com/codecrafters-io/redis-starter-go/app/internal/storage"
)

var (
	memoryStorage = storage.NewMemoryStorage()
	registry      = commands.NewRegistry()
	serverInfo    *storage.Information
)

func getMasterAddress(replicaof string) (string, string, bool, error) {
	fmt.Println(strings.TrimSpace(replicaof))
	address := strings.Split(replicaof, " ")
	if len(address[0]) == 0 {
		return "", "", false, nil
	}
	if len(address) != 2 {
		return "", "", false, fmt.Errorf("replicaof must follow this  pattern localhost:12345")
	}

	return address[0], address[1], true, nil
}

func startReplicationHandshake(conn net.Conn) error {
	_, err := conn.Write([]byte("*1\r\n$4\r\nPING\r\n"))
	if err != nil {
		fmt.Println("Error writing:", err)
		return err
	}
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err)
		return err
	}
	response := string(buf[:n])
	fmt.Println("replication ping response:", response)
	return nil
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	port := flag.Int("port", 6379, "Port for TCP server; 6379 as default")
	masterInfo := flag.String("replicaof", "", "Port for Redis Master")
	flag.Parse()

	masterAddress, masterPort, isSlave, err := getMasterAddress(*masterInfo)
	if err != nil {
		fmt.Println(err)
		return
	}

	var connectionWithMaster string
	connection := "0.0.0.0:" + strconv.Itoa(*port)
	if isSlave {
		serverInfo = storage.NewInformation("slave", "9991b4fb1155b71f4a04d3e1bc3e18c4a990aeeb", "0")
		connectionWithMaster = masterAddress + ":" + masterPort
		fmt.Println("Sending ping to master", connectionWithMaster)

		master, _ := net.Dial("tcp", connectionWithMaster)
		defer master.Close()

		startReplicationHandshake(master)
	} else {
		serverInfo = storage.NewInformation("master", "8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb", "0")
	}

	listener, err := net.Listen("tcp", connection)
	if err != nil {
		fmt.Println("Failed to bind to port ", port)
		return
	}

	defer listener.Close()

	fmt.Printf("Server is listening in port: %d\n", *port)

	registry.Register(&commands.PingCommand{})
	registry.Register(&commands.EchoCommand{})
	registry.Register(&commands.GetCommand{Storage: memoryStorage})
	registry.Register(&commands.SetCommand{Storage: memoryStorage})
	registry.Register(&commands.InfoCommand{Info: serverInfo})

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

		cmd.Execute(conn)
	}
}

func parseCommand(request string) (commands.Command, error) {
	messages := strings.Split(request, CRLF)

	cmdString := strings.ToLower(messages[2])

	cmd, err := registry.Get(cmdString)
	if err != nil {
		return nil, err
	}

	err = cmd.ParseArguments(messages)
	if err != nil {
		return nil, err
	}

	return cmd, nil
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
