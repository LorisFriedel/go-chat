package console

import (
	"github.com/LorisFriedel/go-chat/core"

	log "github.com/sirupsen/logrus"
)

func init() {
	registerCmd("forget", newCmdForget)
}

func newCmdForget(client core.IClient, provider IProvider) (Command, error) {
	chanName, err := provider.GetString()
	if err != nil {
		log.Errorln("newCmdForget: can't get 'chanName' args for instantiating command")
		return nil, err
	}

	return func() error {
		return client.Forget(chanName)
	}, nil
}
