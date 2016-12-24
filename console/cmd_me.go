package console

import (
	"github.com/LorisFriedel/go-chat/core"
)

func init() {
	registerCmd("me", newCmdMe)
}

func newCmdMe(client *core.Client, provider IProvider) (Command, error) {
	return func() error {
		return client.Me()
	}, nil
}
