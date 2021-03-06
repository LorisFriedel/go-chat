package console

import (
	"github.com/LorisFriedel/go-chat/core"

	log "github.com/sirupsen/logrus"
)

func init() {
	registerCmd("close", newCmdClose)
}

func newCmdClose(client core.IClient, provider IProvider) (Command, error) {
	chanName, err := provider.GetString()
	if err != nil {
		log.Errorln("newCmdClose: can't get 'chanName' args for instantiating command")
		return nil, err
	}

	return func() error {
		return client.CloseChan(chanName)
	}, nil
}
