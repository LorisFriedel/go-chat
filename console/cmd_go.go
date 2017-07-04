package console

import (
	"github.com/LorisFriedel/go-chat/core"

	log "github.com/sirupsen/logrus"
)

func init() {
	registerCmd("go", newCmdGo)
}

func newCmdGo(client *core.Client, provider IProvider) (Command, error) {
	name, err := provider.GetString()
	if err != nil {
		log.Errorln("newCmdGo: can't get 'name' args for instantiating command")
		return nil, err
	}

	if provider.HasMore() {
		address, err := provider.GetString()
		if err != nil {
			log.Errorln("newCmdGo: can't get 'address' args for instantiating command")
			return nil, err
		}

		port, err := provider.GetInt()
		if err != nil {
			log.Errorln("newCmdGo: can't get 'port' args for instantiating command")
			return nil, err
		}

		password, err := provider.GetString()
		if err != nil {
			log.Errorln("newCmdGo: can't get 'password' args for instantiating command")
			return nil, err
		}

		return func() error {
			return client.Connect(name, address, port, password)
		}, nil
	}

	return func() error {
		return client.ConnectKnown(name)
	}, nil
}
