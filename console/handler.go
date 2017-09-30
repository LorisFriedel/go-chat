package console

import (
	"fmt"
	"github.com/LorisFriedel/go-chat/core"

	log "github.com/sirupsen/logrus"
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
		log.Errorf("Handler.handle: error while parsing input: %s\n", input)
		return err
	}

	// Find command name
	cmdName := provider.CommandName()

	// Create executable command
	cmd, err := NewCommand(client, cmdName, provider)
	if err != nil {
		log.Errorf("Handler.handle: error while instanciating command: %s\n", cmdName)
		return err
	}

	// Command execution
	if cmd == nil {
		return fmt.Errorf("command %s can't be nil", cmdName)
	}

	err = cmd()

	if err != nil {
		// TODO special case for suicide error, that is not really an error
		log.Errorf("Handler.handle: error while executing command: %s\n", cmdName)
		return err
	}

	return nil // no error, successful handling
}
