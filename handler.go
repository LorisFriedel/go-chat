package main

type IHandler interface {
	Handle(input string) error
}

type ChatHandler struct {

}

func (h *ChatHandler) Handle(input string) error {
	// parse
	// create command
	// execute it
	return nil
}
// TODO
