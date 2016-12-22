package console

import (
	"errors"
	"fmt"
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

type commandFactory func(client *core.Client, provider IProvider) (ICommand, error)

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
		log.Panicf("Command factory %s does not exist.", name)
	}
	m := CommandFactories()
	_, registered := m[name]
	if registered {
		fmt.Errorf("Command factory %s already registered. Replacing.", name)
	}
	m[name] = factory
}

func init() {
	register("message", NewCmdMessage)
	//register("cmd1", NewMemoryDataStore)
}

func NewCommand(client *core.Client, name string, provider IProvider) (ICommand, error) {
	cmdFactory, set := FactoryFor(name)
	if !set {
		// Factory has not been registered
		return nil, errors.New("Invalid Command name: " + name)
	}

	// Run the factory with args
	return cmdFactory(client, provider)
}

/*
Quand on veut handle, on aimerai que la commande soit créer toutes seules (pattern factory ?)
On lui donne un provider, il l'utilise et nous renvoi une erreur
si lors de son extraction de donnée du provider il y a une erreur
example:

provider = newprovider("create Channel1 Password1") // on pourra faire trois provider.GetString()
cmd := NewCommand("chan", provider)

*/
