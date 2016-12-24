package console

import (
	"fmt"
	"github.com/LorisFriedel/go-chat/core"
	"github.com/golang/glog"
	"log"
	"sync"
)

type ICommand interface {
	Execute() error
}

type Cmd struct {
	client *core.Client
}

type commandFactory func(client *core.Client, provider IProvider) (ICommand, error)

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
		glog.Infof("command factory %s already registered. Replacing.\n", name)
	}
	m[name] = factory
}

func NewCommand(client *core.Client, name string, provider IProvider) (ICommand, error) {
	cmdFactory, set := CmdFactory(name)
	if !set {
		// Factory has not been registered
		glog.Errorf("NewCommand: command factory not available for command: %s\n", name)

		return nil, fmt.Errorf("invalid command name: %s", name)
	}

	// Run the factory with args
	return cmdFactory(client, provider)
}
