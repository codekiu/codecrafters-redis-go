package main

import (
	"fmt"
	"strconv"
	"strings"

	// Uncomment this block to pass the first stage
	"net"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")
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
			return
		}

		request := strings.ToLower(strings.TrimSpace(string(buf[:n])))
		fmt.Println("Request is: ", request)

		str := "+PONG\r\n"
		response := []byte(str)

		conn.Write(response)
	}
}
