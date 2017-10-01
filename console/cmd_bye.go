package console

import (
	"github.com/LorisFriedel/go-chat/core"
)

func init() {
	registerCmd("bye", newCmdBye)
}

func newCmdBye(client core.IClient, provider IProvider) (Command, error) {
	return func() error {
		return client.Bye()
	}, nil
}
