package console
import (
	"github.com/LorisFriedel/go-chat/core"
	"github.com/golang/glog"
)

func init() {
	registerCmd("close", newCmdClose)
}

func newCmdClose(client *core.Client, provider IProvider) (Command, error) {
	chanName, err := provider.GetString()
	if err != nil {
		glog.Errorln("newCmdClose: can't get 'chanName' args for instantiating command")
		return nil, err
	}

	return func() error {
		return client.CloseChan(chanName)
	}, nil
}

