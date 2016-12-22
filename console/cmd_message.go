package console

import (
	"github.com/LorisFriedel/go-chat/core"
)

type CmdMessage struct {
	Cmd
	message string
	channel *core.Channel
}

func NewCmdMessage(client *core.Client, provider IProvider) (ICommand, error) {
	// TODO pattern ? -> TO FACTOR
	msg, err := provider.NextString()
	if err != nil {
		return nil, err
	}

	// For now, user can't send message to a channel he is connected
	//channel, err := provider.NextString()
	//if err != nil {
	//	return nil, err
	//}

	return &CmdMessage{Cmd{client}, msg, client.Channel}, nil
}

func (c *CmdMessage) Execute() error {
	c.client.SendMessage(c.message, c.channel) // TODO
	return nil // TODO
}
