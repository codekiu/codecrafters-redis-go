package command

import (
	"fmt"
	"net"
	"sync"

	"github.com/codecrafters-io/redis-starter-go/app/internal/storage"
)

type Command interface {
	Execute(conn net.Conn) error

	Name() string

	ParseArguments(args []string) error
}

type Registry struct {
	commands map[string]Command
	mx       sync.RWMutex
}

func NewRegistry() *Registry {
	return &Registry{
		commands: make(map[string]Command),
		mx:       sync.RWMutex{},
	}
}

func (r *Registry) RegisterCommands(storage *storage.MemoryStorage, info *storage.Information) {
	r.Register(NewPingCommand())
	r.Register(NewEchoCommand())
	r.Register(NewReplconfCommand())
	r.Register(NewInfoCommand(info))
	r.Register(NewSetCommand(storage))
	r.Register(NewGetCommand(storage))
}

func (r *Registry) Register(command Command) error {
	if _, ok := r.commands[command.Name()]; ok {
		return fmt.Errorf("Command %s already registered", command.Name())
	}
	r.commands[command.Name()] = command
	return nil
}

func (r *Registry) Get(name string) (Command, error) {
	cmd, ok := r.commands[name]
	if !ok {
		return nil, fmt.Errorf("command not found: %s", name)
	}
	return cmd, nil
}
