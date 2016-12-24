package console

import (
	"github.com/LorisFriedel/go-chat/core"
	"github.com/golang/glog"
)

type IHandler interface {
	Handle(client *core.Client, input string) error
}

type CmdHandler struct {
	parser IParser
}

func NewHandler(parser IParser) *CmdHandler {
	return &CmdHandler{parser}
}

func (h *CmdHandler) Handle(client *core.Client, input string) error {
	// Parse input
	provider, err := h.parser.Parse(input)
	if err != nil {
		glog.Errorf("Handler.Handle: error while parsing input: %s\n", input)
		return err
	}

	// Find command name
	cmdName := provider.CommandName()

	// Create executable command
	cmd, err := NewCommand(client, cmdName, provider)
	if err != nil {
		glog.Errorf("Handler.Handle: error while instanciating command: %s\n", cmdName)
		return err
	}

	// Execute command
	err = cmd.Execute()
	if err != nil {
		glog.Errorf("Handler.Handle: error while executing command: %s\n", cmdName)
		return err
	}

	return nil // no error, successful handling
}
