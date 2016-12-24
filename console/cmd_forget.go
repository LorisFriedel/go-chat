package console

import (
	"github.com/LorisFriedel/go-chat/core"
	"github.com/golang/glog"
)

func init() {
	registerCmd("forget", newCmdForget)
}

func newCmdForget(client *core.Client, provider IProvider) (Command, error) {
	chanName, err := provider.GetString()
	if err != nil {
		glog.Errorln("newCmdForget: can't get 'chanName' args for instantiating command")
		return nil, err
	}

	return func() error {
		return client.Forget(chanName)
	}, nil
}
