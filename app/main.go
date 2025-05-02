package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/internal/replication"
	"github.com/codecrafters-io/redis-starter-go/app/internal/server"
	"github.com/codecrafters-io/redis-starter-go/app/internal/storage"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	port := flag.String("port", "6379", "Port to listen on")
	replicaOf := flag.String("replicaof", "", "Replica of another Redis instance (host:port)")
	flag.Parse()

	masterHost, masterPort, isSlave, err := getMasterAddress(*replicaOf)
	if err != nil {
		fmt.Println(err)
		return
	}

	var serverInfo *storage.Information
	memoryStorage := storage.NewMemoryStorage()
	fmt.Println("Init new server")

	if isSlave {
		serverInfo = storage.NewInformation("slave", "whatver", "0")
	} else {
		serverInfo = storage.NewInformation("master", "1", "0")
	}

	server, err := server.NewServer("0.0.0.0", *port, memoryStorage, serverInfo)
	if err != nil {
		fmt.Println(err)
		return
	}

	if isSlave {
		replica := replication.NewReplicationManager(masterHost, masterPort, server)
		replica.StartReplication()
	}

	server.Start()
}

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
