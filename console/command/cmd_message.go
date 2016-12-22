package command

import (
	"github.com/LorisFriedel/go-chat/console/provider"
	"github.com/LorisFriedel/go-chat/core"
)

type CmdMessage struct {
	Cmd
	message string
}

func NewCmdMessage(client *core.Client, provider provider.IProvider) (ICommand, error) {
	// TODO pattern ? -> TO FACTOR
	msg, err := provider.GetString()
	if err != nil {
		return nil, err
	}

	// For now, user can't send message to a channel he is connected
	//channel, err := provider.NextString()
	//if err != nil {
	//	return nil, err
	//}

	return &CmdMessage{Cmd{client}, msg}, nil
}

func (c *CmdMessage) Execute() error {
	c.client.SendMessage(c.message) // TODO
	return nil                      // TODO
}
