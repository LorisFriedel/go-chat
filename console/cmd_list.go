package console

import (
	"github.com/LorisFriedel/go-chat/core"
)

func init() {
	registerCmd("list", newCmdList)
}

func newCmdList(client *core.Client, provider IProvider) (Command, error) {
	return func() error {
		return client.List()
	}, nil
}
