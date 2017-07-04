package console

import (
	"github.com/LorisFriedel/go-chat/core"

	log "github.com/sirupsen/logrus"
)

func init() {
	registerCmd("message", newCmdMessage)
}

func newCmdMessage(client *core.Client, provider IProvider) (Command, error) {
	text, err := provider.GetString()
	if err != nil {
		log.Errorln("newCmdMessage: can't get 'text' args for instantiating command")
		return nil, err
	}

	return func() error {
		return client.SendMessage(text)
	}, nil
}
