package console

import (
	"fmt"
	"sync"

	"github.com/LorisFriedel/go-chat/core"

	log "github.com/sirupsen/logrus"
)

type Command func() error

type commandFactory func(client core.IClient, provider IProvider) (Command, error)

// No need to sync the map, no concurrency there
var commandFactories map[string]commandFactory
var onceCommandFactories sync.Once

func CommandFactories() map[string]commandFactory {
	onceCommandFactories.Do(func() {
		commandFactories = make(map[string]commandFactory)
	})
	return commandFactories
}

func CmdFactory(name string) (commandFactory, bool) {
	m := CommandFactories()
	val, set := m[name]
	return val, set
}

func registerCmd(name string, factory commandFactory) {
	if factory == nil {
		log.Panicf("command factory %s does not exist.\n", name)
	}
	m := CommandFactories()
	_, registered := m[name]
	if registered {
		log.Infof("command factory %s already registered. Replacing.\n", name)
	}
	m[name] = factory
}

func NewCommand(client core.IClient, name string, provider IProvider) (Command, error) {
	cmdFactory, set := CmdFactory(name)
	if !set {
		// Factory has not been registered
		log.Errorf("NewCommand: command factory not available for command: %s\n", name)

		return nil, fmt.Errorf("invalid command name: %s", name)
	}

	// Run the factory with args
	return cmdFactory(client, provider)
}
