package console

import (
	"github.com/LorisFriedel/go-chat/core"
)

func init() {
	registerCmd("me", newCmdMe)
}

func newCmdMe(client core.IClient, provider IProvider) (Command, error) {
	return func() error {
		return client.Me()
	}, nil
}
