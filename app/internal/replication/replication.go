package replication

import (
	"fmt"
	"net"

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

	_, err = conn.Write([]byte("*1\r\n$4\r\nPING\r\n"))
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

	_, err = conn.Write([]byte("*3\r\n$8\r\nREPLCONF\r\n$14\r\nlistening-port\r\n$4\r\n6380\r\n"))
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

	_, err = conn.Write([]byte("*3\r\n$8\r\nREPLCONF\r\n$4\r\ncapa\r\n$6\r\npsync2\r\n"))
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
