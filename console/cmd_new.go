package console

import (
	"github.com/LorisFriedel/go-chat/core"
	"github.com/golang/glog"
)

type CmdNew struct {
	Cmd
	name    string
	address string
	port    int
	passwd  string
}

func init() {
	registerCmd("new", newCmdNew)
}

func newCmdNew(client *core.Client, provider IProvider) (ICommand, error) {
	// TODO Identify proper pattern & factor.

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

	return &CmdNew{
		Cmd:     Cmd{client},
		name:    name,
		address: address,
		port:    port,
		passwd:  passwd,
	}, nil
}

func (c *CmdNew) Execute() error {
	glog.Infoln("CmdNew: executing command")
	return c.client.CreateChan(c.name, c.address, c.port, c.passwd)
}
