package replication

import (
	"fmt"
	"net"
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/app/internal/protocol"
	"github.com/codecrafters-io/redis-starter-go/app/internal/server"
)

type ReplicationManager struct {
	host, port string
	server     *server.Server
}

func NewReplicationManager(host, port string, server *server.Server) *ReplicationManager {
	return &ReplicationManager{
		host:   host,
		port:   port,
		server: server,
	}
}

func (r *ReplicationManager) StartReplication() {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", r.host, r.port))
	if err != nil {
		fmt.Printf("Failed to connect to replica: %v\n", err)
		return
	}
	defer conn.Close()

	command := "PING"
	respArray := protocol.T_ARRAY +
		"1" +
		protocol.CRLF +
		protocol.T_BULK_STRING +
		"4" +
		protocol.CRLF +
		command +
		protocol.CRLF

	_, err = conn.Write([]byte(respArray))
	if err != nil {
		fmt.Println("Error writing:", err)
		return
	}
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err)
		return
	}
	response := string(buf[:n])
	fmt.Println("replication ping response:", response)

	command = "REPLCONF"
	subCommand := "listening-port"
	respArray = protocol.T_ARRAY +
		"3" +
		protocol.CRLF +
		protocol.T_BULK_STRING +
		strconv.Itoa(len(command)) +
		protocol.CRLF +
		command +
		protocol.CRLF +
		protocol.T_BULK_STRING +
		strconv.Itoa(len(subCommand)) +
		protocol.CRLF +
		subCommand +
		protocol.CRLF +
		protocol.T_BULK_STRING +
		strconv.Itoa(len(r.port)) +
		protocol.CRLF +
		r.port +
		protocol.CRLF

	_, err = conn.Write([]byte(respArray))
	if err != nil {
		fmt.Println("Error writing:", err)
		return
	}
	buf = make([]byte, 1024)
	n, err = conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err)
		return
	}
	response = string(buf[:n])
	fmt.Println("replication ping response:", response)

	command = "REPLCONF"
	subCommand = "capa"
	respArray = protocol.T_ARRAY +
		"3" +
		protocol.CRLF +
		protocol.T_BULK_STRING +
		strconv.Itoa(len(command)) +
		protocol.CRLF +
		command +
		protocol.CRLF +
		protocol.T_BULK_STRING +
		strconv.Itoa(len(subCommand)) +
		protocol.CRLF +
		subCommand +
		protocol.CRLF +
		protocol.T_BULK_STRING +
		strconv.Itoa(len("psync2")) +
		protocol.CRLF +
		"psync2" +
		protocol.CRLF

	_, err = conn.Write([]byte(respArray))
	if err != nil {
		fmt.Println("Error writing:", err)
		return
	}
	buf = make([]byte, 1024)
	n, err = conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err)
		return
	}
	response = string(buf[:n])
	fmt.Println("replication ping response:", response)

	command = "PSYNC"
	subCommand = "?"
	respArray = protocol.T_ARRAY +
		"3" +
		protocol.CRLF +
		protocol.T_BULK_STRING +
		strconv.Itoa(len(command)) +
		protocol.CRLF +
		command +
		protocol.CRLF +
		protocol.T_BULK_STRING +
		strconv.Itoa(len(subCommand)) +
		protocol.CRLF +
		subCommand +
		protocol.CRLF +
		protocol.T_BULK_STRING +
		strconv.Itoa(len("-1")) +
		protocol.CRLF +
		"-1" +
		protocol.CRLF

	_, err = conn.Write([]byte(respArray))
	if err != nil {
		fmt.Println("Error writing:", err)
		return
	}
	buf = make([]byte, 1024)
	n, err = conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err)
		return
	}
	response = string(buf[:n])
	fmt.Println("replication ping response:", response)
}
