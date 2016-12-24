package console

import (
	"github.com/LorisFriedel/go-chat/core"
	"github.com/golang/glog"
)

func init() {
	registerCmd("new", newCmdNew)
}

func newCmdNew(client *core.Client, provider IProvider) (Command, error) {
	name, err := provider.GetString()
	if err != nil {
		glog.Errorln("newCmdNew: can't get 'name' args for instantiating command")
		return nil, err
	}

	address, err := provider.GetString()
	if err != nil {
		glog.Errorln("newCmdNew: can't get 'address' args for instantiating command")
		return nil, err
	}

	port, err := provider.GetInt()
	if err != nil {
		glog.Errorln("newCmdNew: can't get 'port' args for instantiating command")
		return nil, err
	}

	passwd, err := provider.GetString()
	if err != nil {
		glog.Errorln("newCmdNew: can't get 'passwd' args for instantiating command")
		return nil, err
	}

	return func() error {
		return client.CreateChan(name, address, port, passwd)
	}, nil
}
