package core

type IChannel interface {
	// TODO
}

type Channel struct {
	Name    string
	Address string
	Port    int

	// Name string
	// Address string
	// Port int
	// Socket ?
	// Members ?
}

func NewChannel(name, address string, port int) *Channel {
	// TODO
	return nil // TODO
}
