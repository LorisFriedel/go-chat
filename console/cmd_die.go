package console

import (
	"github.com/LorisFriedel/go-chat/core"
)

func init() {
	registerCmd("die", newCmdDie)
}

func newCmdDie(client *core.Client, provider IProvider) (Command, error) {
	return func() error {
		return client.Die()
	}, nil
}
