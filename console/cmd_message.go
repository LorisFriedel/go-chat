package console

import (
	"github.com/LorisFriedel/go-chat/core"
	"github.com/golang/glog"
)

func init() {
	registerCmd("message", newCmdMessage)
}

func newCmdMessage(client *core.Client, provider IProvider) (Command, error) {
	text, err := provider.GetString()
	if err != nil {
		glog.Errorln("newCmdMessage: can't get 'text' args for instantiating command")
		return nil, err
	}

	return func() error {
		return client.SendMessage(text)
	}, nil
}
