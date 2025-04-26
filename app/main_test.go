package main

import (
	"net"
	"testing"
	"time"
)

func newMockConn() (net.Conn, net.Conn) {
	return net.Pipe()
}

func readResponse(conn net.Conn) string {
	buf := make([]byte, 1024)
	n, _ := conn.Read(buf)
	return string(buf[:n])
}

func TestPingCommand(t *testing.T) {
	server, client := newMockConn()
	defer server.Close()
	defer client.Close()

	cmd := &pingCommand{}
	go cmd.Handle(server)

	response := readResponse(client)
	expected := T_SIMPLE_STRING + "PONG" + CRLF

	if response != expected {
		t.Errorf("expected %q but got %q", expected, response)
	}
}

func TestEchoCommand(t *testing.T) {
	server, client := newMockConn()
	defer server.Close()
	defer client.Close()

	stringToSend := "Hello!"

	cmd := &echoCommand{Content: stringToSend}
	go cmd.Handle(server)

	response := readResponse(client)
	expected := T_SIMPLE_STRING + stringToSend + CRLF

	if response != expected {
		t.Errorf("expected %q but got %q", expected, response)
	}
}

func TestHandleClientEcho(t *testing.T) {
	server, client := newMockConn()
	defer server.Close()
	message := T_ARRAY + "2" + CRLF +
		T_BULK_STRING + "4" + CRLF + "echo" + CRLF +
		T_BULK_STRING + "6" + CRLF + "Hello!" + CRLF

	// Write the command from the client side
	go func() {
		client.Write([]byte(message))
	}()

	go handleClient(server)

	response := readResponse(client)
	expected := T_SIMPLE_STRING + "Hello!" + CRLF

	if response != expected {
		t.Errorf("expected %q but got %q", expected, response)
	}

	client.Close()
}

func TestHandleClientPing(t *testing.T) {
	server, client := newMockConn()
	defer server.Close()
	defer client.Close()
	message := T_ARRAY + "1" + CRLF +
		T_BULK_STRING + "4" + CRLF + "ping" + CRLF
	go handleClient(server)

	// Write the command from the client side
	client.Write([]byte(message))

	response := readResponse(client)
	expected := T_SIMPLE_STRING + "PONG" + CRLF

	if response != expected {
		t.Errorf("expected %q but got %q", expected, response)
	}
}

func TestHandleClientSet(t *testing.T) {
	resetDict()

	server, client := newMockConn()
	defer server.Close()
	defer client.Close()

	// Define the Redis SET command message
	message := T_ARRAY + "3" + CRLF +
		T_BULK_STRING + "3" + CRLF + "set" + CRLF +
		T_BULK_STRING + "5" + CRLF + "hello" + CRLF +
		T_BULK_STRING + "4" + CRLF + "hola" + CRLF

	// Start handling in a goroutine
	go func() {
		// Add a small delay to ensure goroutine starts after client writes
		time.Sleep(10 * time.Millisecond)
		handleClient(server)
	}()

	// Write the command from the client side
	_, err := client.Write([]byte(message))
	if err != nil {
		t.Fatalf("Failed to write to client: %v", err)
	}

	// Read the response from the server
	response := readResponse(client)
	expected := T_SIMPLE_STRING + "OK" + CRLF

	// Validate the response
	if response != expected {
		t.Errorf("expected %q but got %q", expected, response)
	}
}

func TestHandleClientGetEmpty(t *testing.T) {
	resetDict()

	server, client := newMockConn()
	defer server.Close()
	defer client.Close()

	// Define the Redis SET command message
	message := T_ARRAY + "3" + CRLF +
		T_BULK_STRING + "3" + CRLF + "get" + CRLF +
		T_BULK_STRING + "5" + CRLF + "hello" + CRLF

	// Start handling in a goroutine
	go func() {
		// Add a small delay to ensure goroutine starts after client writes
		time.Sleep(10 * time.Millisecond)
		handleClient(server)
	}()

	// Write the command from the client side
	_, err := client.Write([]byte(message))
	if err != nil {
		t.Fatalf("Failed to write to client: %v", err)
	}

	// Read the response from the server
	response := readResponse(client)
	expected := T_BULK_STRING + "-1" + CRLF

	// Validate the response
	if response != expected {
		t.Errorf("expected %q but got %q", expected, response)
	}
}

func TestHandleClientGet(t *testing.T) {
	resetDict()

	key := "hello"
	value := "hola"

	dict[key] = value

	server, client := newMockConn()
	defer server.Close()
	defer client.Close()

	// Define the Redis SET command message
	message := T_ARRAY + "3" + CRLF +
		T_BULK_STRING + "3" + CRLF + "get" + CRLF +
		T_BULK_STRING + "5" + CRLF + key + CRLF

	// Start handling in a goroutine
	go func() {
		// Add a small delay to ensure goroutine starts after client writes
		time.Sleep(10 * time.Millisecond)
		handleClient(server)
	}()

	// Write the command from the client side
	_, err := client.Write([]byte(message))
	if err != nil {
		t.Fatalf("Failed to write to client: %v", err)
	}

	// Read the response from the server
	response := readResponse(client)
	expected := T_SIMPLE_STRING + value + CRLF

	// Validate the response
	if response != expected {
		t.Errorf("expected %q but got %q", expected, response)
	}
}

func resetDict() {
	for k := range dict {
		delete(dict, k)
	}
}
