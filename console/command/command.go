package command

import (
	"fmt"
	"github.com/LorisFriedel/go-chat/console/provider"
	"github.com/LorisFriedel/go-chat/core"
	"log"
	"sync"
)

type ICommand interface {
	Execute() error
}

type Cmd struct {
	client *core.Client
}

type commandFactory func(client *core.Client, provider provider.IProvider) (ICommand, error)

// No need to sync the map, no concurrency there
var commandFactories map[string]commandFactory
var once sync.Once

func CommandFactories() map[string]commandFactory {
	once.Do(func() {
		commandFactories = make(map[string]commandFactory)
	})
	return commandFactories
}

func FactoryFor(name string) (commandFactory, bool) {
	m := CommandFactories()
	val, set := m[name]
	return val, set
}

func register(name string, factory commandFactory) {
	if factory == nil {
		log.Panicf("command factory %s does not exist.", name)
	}
	m := CommandFactories()
	_, registered := m[name]
	if registered {
		log.Printf("command factory %s already registered. Replacing.", name)
	}
	m[name] = factory
}

func init() {
	register("message", NewCmdMessage)
	//register("cmd1", NewMemoryDataStore)
}

func New(client *core.Client, name string, provider provider.IProvider) (ICommand, error) {
	cmdFactory, set := FactoryFor(name)
	if !set {
		// Factory has not been registered
		return nil, fmt.Errorf("command.New: invalid command name: %s", name)
	}

	// Run the factory with args
	return cmdFactory(client, provider)
}
