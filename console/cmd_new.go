package console

import (
	"github.com/LorisFriedel/go-chat/core"

	log "github.com/sirupsen/logrus"
)

func init() {
	registerCmd("new", newCmdNew)
}

func newCmdNew(client core.IClient, provider IProvider) (Command, error) {
	name, err := provider.GetString()
	if err != nil {
		log.Errorln("newCmdNew: can't get 'name' args for instantiating command")
		return nil, err
	}

	address, err := provider.GetString()
	if err != nil {
		log.Errorln("newCmdNew: can't get 'address' args for instantiating command")
		return nil, err
	}

	port, err := provider.GetInt()
	if err != nil {
		log.Errorln("newCmdNew: can't get 'port' args for instantiating command")
		return nil, err
	}

	var password string
	if provider.HasMore() {
		password, err = provider.GetString()
		if err != nil {
			log.Errorln("newCmdNew: can't get 'password' args for instantiating command")
			return nil, err
		}
	}

	var timeout int
	if provider.HasMore() {
		timeout, err = provider.GetInt()
		if err != nil {
			log.Errorln("newCmdNew: can't get 'timeout' args for instantiating command")
			return nil, err
		}
	}

	return func() error {
		return client.CreateConnectChan(name, address, port, password, timeout)
	}, nil
}
