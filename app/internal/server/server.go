package server

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/internal/command"
	"github.com/codecrafters-io/redis-starter-go/app/internal/protocol"
	"github.com/codecrafters-io/redis-starter-go/app/internal/storage"
)

type Server struct {
	listener net.Listener
	storage  storage.Storage
	// replica *ReplicationManager
	commandReg *command.Registry
}

func NewServer(host, port string, storage *storage.MemoryStorage, info *storage.Information) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", host, port))
	fmt.Println("new server listening")
	if err != nil {
		fmt.Println("failed")
		return nil, err
	}

	reg := command.NewRegistry()
	fmt.Println("registering commands")
	reg.RegisterCommands(storage, info)

	fmt.Println("returning")
	return &Server{
		listener:   listener,
		storage:    storage,
		commandReg: reg,
	}, nil
}

func (s *Server) Start() {
	log.Printf("Server listening on port %s", s.listener.Addr())
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Printf("Accept error %v:", err)
		}

		go s.handleConnection(conn)
	}
}

func (s *Server) GetPort() string {
	addr := strings.Split(s.listener.Addr().String(), ":")
	return addr[len(addr)-1]
}

func (s *Server) handleConnection(conn net.Conn) {
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

		cmd, innerErr := parseCommand(s, string(buf[:n]))
		if innerErr != nil {
			conn.Write([]byte(protocol.T_SIMPLE_STRING + innerErr.Error() + protocol.CRLF))
			continue
		}

		cmd.Execute(conn)
	}
}

func parseCommand(server *Server, request string) (command.Command, error) {
	messages := strings.Split(request, protocol.CRLF)

	cmdString := strings.ToLower(messages[2])

	cmd, err := server.commandReg.Get(cmdString)
	if err != nil {
		return nil, err
	}

	err = cmd.ParseArguments(messages)
	if err != nil {
		return nil, err
	}

	return cmd, nil
}
