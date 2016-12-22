package console

import (
	"github.com/LorisFriedel/go-chat/core"
	"fmt"
)

type IHandler interface {
	Handle(client *core.Client, input string) error
}

type CmdHandler struct {
	parser IParser
}

func NewCmdHandler(parser IParser) *CmdHandler {
	return &CmdHandler{parser}
}

func (h *CmdHandler) Handle(client *core.Client, input string) error {
	// Parse input
	provider, err := h.parser.Parse(input)
	if err != nil {
		// TODO handle error
	}

	// Find command name
	cmdName := provider.CommandName()

	fmt.Println("Commande name : ", cmdName) // TODO remove

	// Create executable command
	command, err := NewCommand(client, cmdName, provider)
	if err != nil {
		// TODO handle error
	}

	// Execute command
	return command.Execute()
}
