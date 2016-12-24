package console

import (
	"github.com/LorisFriedel/go-chat/core"
	"github.com/golang/glog"
)

func init() {
	registerCmd("go", newCmdGo)
}

func newCmdGo(client *core.Client, provider IProvider) (Command, error) {
	name, err := provider.GetString()
	if err != nil {
		glog.Errorln("newCmdGo: can't get 'name' args for instantiating command")
		return nil, err
	}

	if provider.HasMore() {
		address, err := provider.GetString()
		if err != nil {
			glog.Errorln("newCmdGo: can't get 'address' args for instantiating command")
			return nil, err
		}

		port, err := provider.GetInt()
		if err != nil {
			glog.Errorln("newCmdGo: can't get 'port' args for instantiating command")
			return nil, err
		}

		return func() error {
			return client.Connect(name, address, port)
		}, nil
	}

	return func() error {
		return client.ConnectKnown(name)
	}, nil
}
