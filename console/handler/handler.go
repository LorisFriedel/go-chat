package handler

import (
	"github.com/LorisFriedel/go-chat/console/command"
	"github.com/LorisFriedel/go-chat/console/parser"
	"github.com/LorisFriedel/go-chat/core"
)

type IHandler interface {
	Handle(client *core.Client, input string) error
}

type CmdHandler struct {
	parser parser.IParser
}

func New(parser parser.IParser) *CmdHandler {
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

	// Create executable command
	cmd, err := command.New(client, cmdName, provider)
	if err != nil {
		return err
	}

	// Execute command
	return cmd.Execute()
}
