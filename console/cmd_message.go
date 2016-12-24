package console

import (
	"github.com/LorisFriedel/go-chat/core"
	"github.com/golang/glog"
)

type CmdMessage struct {
	Cmd
	text string
}

func init() {
	registerCmd("message", newCmdMessage)
}

func newCmdMessage(client *core.Client, provider IProvider) (ICommand, error) {
	// TODO pattern ? -> TO FACTOR
	text, err := provider.GetString()
	if err != nil {
		glog.Errorln("newCmdMessage: can't get 'text' args for instantiating command")
		return nil, err
	}

	return &CmdMessage{
		Cmd:  Cmd{client},
		text: text,
	}, nil
}

func (c *CmdMessage) Execute() error {
	glog.Infoln("CmdMessage: executing command")
	return c.client.SendMessage(c.text)
}
